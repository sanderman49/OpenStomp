using System;

namespace OpenStomp.Models.Pedal;

public class Led
{
    private bool _state;
    public bool State { get => _state; }

    private LedMode _mode;
    public LedMode Mode { get => _mode; }
    
    public enum LedMode
    {
        Normal = 0,
        Bpm = 1,
    }

    public Led()
    {
        _mode = LedMode.Normal;
    }

    public void SetState(bool state)
    {
        if (Mode == LedMode.Normal)
        {
            _state = state;
        }
    }

    public void SetBpm(float bpm)
    {
        if (Mode == LedMode.Bpm)
        {
            throw new NotImplementedException();
        }
    }

    public void UpdateMode(LedMode mode)
    {
         _mode = mode;
    }
}