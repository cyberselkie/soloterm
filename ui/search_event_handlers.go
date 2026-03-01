package ui

func (a *App) handleSearchCancelled(e *SearchCancelledEvent) {
	a.pages.HidePage(SEARCH_MODAL_ID)
	a.SetFocus(a.sessionView.TextArea)
}

func (a *App) handleSearchShow(e *SearchShowEvent) {
	a.searchView.Reset()
	a.pages.ShowPage(SEARCH_MODAL_ID)
	a.SetFocus(a.searchView.searchTermInput)
}

func (a *App) handleSearchSelectResult(e *SearchSelectResultEvent) {
	match := a.searchView.CurrentMatch()
	if match == nil {
		return
	}
	term := a.searchView.lastTerm

	a.pages.HidePage(SEARCH_MODAL_ID)

	if match.isNotes {
		g := a.CurrentGame()
		if g == nil {
			return
		}
		a.HandleEvent(&GameNotesSelectedEvent{
			BaseEvent: BaseEvent{action: GAME_NOTES_SELECTED},
			GameID:    g.ID,
		})
		a.gameView.SelectNotes(g.ID)
	} else {
		// Load the matched session and highlight it in the tree
		a.sessionView.SelectSession(match.sessionID)
		a.gameView.SelectSession(match.sessionID)
	}

	a.SetFocus(a.sessionView.TextArea)

	// Defer Select to after the TextArea has rendered the new content.
	// Calling Select immediately after SetText uses stale layout and misplaces
	// the cursor (especially on the last line). QueueUpdateDraw must be called
	// from a goroutine — calling it from the main event goroutine deadlocks.
	//
	// This was a big work around to get it to function properly
	offset := match.offset
	go a.QueueUpdateDraw(func() {
		ta := a.sessionView.TextArea
		// Use SetMovedFunc as a one-shot hook: Select calls moved() synchronously
		// after updating cursor.row, so GetCursor() is reliable at that moment.
		ta.SetMovedFunc(func() {
			ta.SetMovedFunc(nil) // one-shot: clear immediately
			fromRow, _, _, _ := ta.GetCursor()
			ta.SetOffset(fromRow, 0)
		})
		ta.Select(offset, offset+len(term))
		// Re-apply the selection in the next draw cycle. When switching sessions
		// SetText resets lineStarts; if reset() fires in this draw (e.g. due to
		// a width change after the modal hides), findCursor collapses
		// selectionStart=cursor. A second Select re-establishes it after the
		// layout has stabilised. The scroll from SetOffset survives either way.
		go a.QueueUpdateDraw(func() {
			ta.Select(offset, offset+len(term))
		})
	})
}
