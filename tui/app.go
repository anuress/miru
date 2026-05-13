package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.design/x/clipboard"

	"github.com/anuress/miru/adb"
	"github.com/anuress/miru/model"
	"github.com/anuress/miru/protocol"
)

type focusPane int

const (
	focusList focusPane = iota
	focusDetail
)

type msgReceived struct{ msg protocol.Message }
type connLost struct{}
type connRestored struct{ session *adb.LogcatSession }

type clearCurlMsg struct{}

type AppModel struct {
	list         ListModel
	detail       DetailModel
	filter       Filter
	filterMode   bool
	filterInput  string
	curlFlash    string // brief status after copying curl
	focus        focusPane
	session      *adb.LogcatSession
	msgCh        <-chan protocol.Message
	device       string
	process      string
	pid          string
	connected    bool
	everReceived bool
	width        int
	height       int
}

func NewAppModel(session *adb.LogcatSession, device, process, pid string) AppModel {
	return AppModel{
		list:    NewListModel(),
		detail:  NewDetailModel(),
		session: session,
		msgCh:     protocol.NewStreamReader(session),
		device:    device,
		process:   process,
		pid:       pid,
		connected: true,
	}
}

func (m AppModel) Init() tea.Cmd {
	return m.listenCmd()
}

func (m AppModel) listenCmd() tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-m.msgCh
		if !ok {
			return connLost{}
		}
		return msgReceived{msg: msg}
	}
}

func (m AppModel) reconnectCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(_ time.Time) tea.Msg {
		session, err := adb.StartLogcat(m.device, m.pid)
		if err != nil {
			return connLost{}
		}
		return connRestored{session: session}
	})
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.list.width = (m.width*45/100) - 4  // border(2) + padding(1) + safety(1)
		m.list.height = m.height - 4
		m.detail.width = m.width - (m.width*45/100) - 4
		m.detail.height = m.height - 4
		return m, nil

	case msgReceived:
		m.everReceived = true
		m.connected = true
		m.applyMessage(msg.msg)
		return m, m.listenCmd()

	case connLost:
		m.connected = false
		if m.session != nil {
			m.session.Stop()
		}
		return m, m.reconnectCmd()

	case connRestored:
		m.session = msg.session
		m.msgCh = protocol.NewStreamReader(msg.session)
		m.connected = true
		return m, m.listenCmd()

	case clearCurlMsg:
		m.curlFlash = ""
		return m, nil

	case tea.KeyMsg:
		if m.filterMode {
			switch msg.String() {
			case "esc":
				m.filterMode = false
				m.filterInput = ""
				m.filter = NewFilter("")
				m.list.SetFilter(m.filter)
			case "backspace":
				if len(m.filterInput) > 0 {
					m.filterInput = m.filterInput[:len(m.filterInput)-1]
					m.filter = NewFilter(m.filterInput)
					m.list.SetFilter(m.filter)
				}
			default:
				if len(msg.String()) == 1 {
					m.filterInput += msg.String()
					m.filter = NewFilter(m.filterInput)
					m.list.SetFilter(m.filter)
				}
			}
			return m, nil
		}
		switch msg.String() {
		case "q", "ctrl+c":
			if m.session != nil {
				m.session.Stop()
			}
			return m, tea.Quit
		case "tab":
			if m.focus == focusList {
				m.focus = focusDetail
			} else {
				m.focus = focusList
			}
		case "f":
			m.filterMode = true
		case "c":
			m.list.Clear()
			m.detail.SetRequest(nil)
		case "y":
			if m.focus == focusDetail {
				if val, ok := m.detail.ValueAtCursor(); ok {
					clipboard.Write(clipboard.FmtText, []byte(val))
					m.curlFlash = "✓ copied"
					return m, tea.Tick(2*time.Second, func(_ time.Time) tea.Msg {
						return clearCurlMsg{}
					})
				}
			} else if sel := m.list.Selected(); sel != nil {
				clipboard.Write(clipboard.FmtText, []byte(GenerateCurl(*sel)))
				m.curlFlash = "✓ curl copied"
				return m, tea.Tick(2*time.Second, func(_ time.Time) tea.Msg {
					return clearCurlMsg{}
				})
			}
		case "Y":
			if m.focus == focusDetail {
				if line, ok := m.detail.LineAtCursor(); ok {
					clipboard.Write(clipboard.FmtText, []byte(line))
					m.curlFlash = "✓ copied"
					return m, tea.Tick(2*time.Second, func(_ time.Time) tea.Msg {
						return clearCurlMsg{}
					})
				}
			}
		default:
			if m.focus == focusList {
				m.list, _ = m.list.Update(msg)
				if sel := m.list.Selected(); sel != nil {
					m.detail.SetRequest(sel)
				}
			} else {
				m.detail, _ = m.detail.Update(msg)
			}
		}
	}
	return m, nil
}

func (m *AppModel) applyMessage(msg protocol.Message) {
	// Skip incomplete messages — can happen from buffered logcat before miru started
	if msg.Method == "" && msg.URL == "" {
		return
	}
	r := model.Request{
		ID:          msg.ID,
		Method:      msg.Method,
		URL:         msg.URL,
		StatusCode:  msg.Status,
		Duration:    time.Duration(msg.Duration) * time.Millisecond,
		ReqHeaders:  msg.ReqHeaders,
		ReqBody:     msg.ReqBody,
		RespHeaders: msg.RespHeaders,
		RespBody:    msg.RespBody,
		InFlight:    false,
	}
	if msg.Error != "" {
		r.Error = msg.Error
	}
	m.list.AddRequest(r)
	if sel := m.list.Selected(); sel != nil {
		m.detail.SetRequest(sel)
	}
}

func (m AppModel) View() string {
	if m.width == 0 {
		return ""
	}

	connStatus := lipgloss.NewStyle().Foreground(ColorGreen).Render("● connected")
	if !m.connected {
		if m.everReceived {
			connStatus = lipgloss.NewStyle().Foreground(ColorRed).Render("● disconnected — retrying…")
		} else {
			connStatus = lipgloss.NewStyle().Foreground(ColorOrange).Render("● waiting for interceptor…")
		}
	}

	filterStr := ""
	if m.filterMode || m.filterInput != "" {
		filterStr = fmt.Sprintf("[Filter: %s_]  ", m.filterInput)
	}

	topBar := lipgloss.NewStyle().Background(ColorBgAlt).Width(m.width).Render(
		fmt.Sprintf(" ◆ miru  %s │ %s  %sf:filter c:clear q:quit", m.process, m.device, filterStr),
	)

	paneH := m.height - 2 // top bar + status bar

	// Each box: border(2) + paddingLeft(1) = 3 overhead on content width
	listInner := (m.width*45/100) - 3
	detailInner := m.width - (m.width*45/100) - 3
	innerH := paneH - 2

	listBorderColor := ColorBorder
	detailBorderColor := ColorBorder
	if m.focus == focusDetail {
		detailBorderColor = ColorBlue
	} else {
		listBorderColor = ColorBlue
	}

	listView := m.list.View()
	detailView := m.detail.View()

	listBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(listBorderColor).
		PaddingLeft(1).
		Width(listInner).Height(innerH).
		Render(listView)

	detailBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(detailBorderColor).
		PaddingLeft(1).
		Width(detailInner).Height(innerH).
		Render(detailView)

	split := lipgloss.JoinHorizontal(lipgloss.Top, listBox, detailBox)

	reqCount := fmt.Sprintf("%d requests", len(m.list.requests))

	right := connStatus
	if m.curlFlash != "" {
		right = lipgloss.NewStyle().Foreground(ColorGreen).Render(m.curlFlash)
	}

	statusBar := lipgloss.NewStyle().Background(ColorBgAlt).Foreground(ColorGray).Width(m.width).Render(
		fmt.Sprintf(" %s │ Tab:panes │ ←→:tabs │ y:copy │ f:filter │ c:clear  %s", reqCount, right),
	)

	return strings.Join([]string{topBar, split, statusBar}, "\n")
}

