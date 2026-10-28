package ipc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"
)

// CheckPlayback changes local playback briefly and leaves it paused.
func (c *Client) CheckPlayback(ctx context.Context, uri string, out io.Writer) error {
	initial, err := c.State(ctx)
	if err != nil {
		return err
	}
	for n := 0; (!initial.Connected || !initial.LoggedIn) && n < 100; n++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
		initial, err = c.State(ctx)
		if err != nil {
			return err
		}
	}
	if !initial.Connected || !initial.LoggedIn {
		return errors.New("Soloist did not restore a paired session; select Spotify TUI Local in an official Spotify client")
	}
	wait := func(predicate func(State) bool) (State, error) {
		deadline := time.NewTimer(15 * time.Second)
		defer deadline.Stop()
		tick := time.NewTicker(100 * time.Millisecond)
		defer tick.Stop()
		for {
			state, e := c.State(ctx)
			if e != nil {
				return state, e
			}
			if predicate(state) {
				return state, nil
			}
			select {
			case <-ctx.Done():
				return state, ctx.Err()
			case <-deadline.C:
				return state, errors.New("provider did not confirm the requested playback state")
			case <-tick.C:
			}
		}
	}
	if err = c.Call(ctx, "play_track", map[string]any{"uri": uri}, nil); err != nil {
		return err
	}
	defer func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = c.Call(stopCtx, "pause", nil, nil)
	}()
	playing, err := wait(func(s State) bool { return s.Status == "playing" && s.URI == uri && s.DurationMS > 0 })
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "Full-track playback state confirmed. Checking progress...")
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
	}
	progress, err := c.State(ctx)
	if err != nil {
		return err
	}
	position := progress.PositionMS
	if progress.Status == "playing" && progress.PositionUpdatedAt > 0 {
		position += time.Now().UnixMilli() - progress.PositionUpdatedAt
	}
	if position <= playing.PositionMS+1000 {
		return errors.New("playback clock did not advance")
	}
	fmt.Fprintln(out, "Playback progress: OK")
	if err = c.Call(ctx, "pause", nil, nil); err != nil {
		return err
	}
	if _, err = wait(func(s State) bool { return s.Status == "paused" }); err != nil {
		return err
	}
	fmt.Fprintln(out, "Pause confirmation: OK")
	if err = c.Call(ctx, "resume", nil, nil); err != nil {
		return err
	}
	if _, err = wait(func(s State) bool { return s.Status == "playing" }); err != nil {
		return err
	}
	fmt.Fprintln(out, "Resume confirmation: OK")
	target := min(int64(15000), playing.DurationMS/2)
	if err = c.Call(ctx, "seek", map[string]any{"position_ms": target}, nil); err != nil {
		return err
	}
	if _, err = wait(func(s State) bool { return s.PositionMS >= target-2000 && s.PositionMS <= target+4000 }); err != nil {
		return err
	}
	fmt.Fprintln(out, "Seek confirmation: OK")
	if err = c.Call(ctx, "set_volume", map[string]any{"volume": initial.Volume}, nil); err != nil {
		return err
	}
	fmt.Fprintln(out, "Volume command accepted at the existing level.")
	return nil
}
