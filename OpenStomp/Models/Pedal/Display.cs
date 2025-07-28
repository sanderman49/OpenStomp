using System;

namespace OpenStomp.Models.Pedal;

public class Display
{
    private string _currentMessage;
    public string CurrentMessage { get; }

    public void DisplayMessage(string text, bool temporary = false)
    {
        _currentMessage = text;
        throw new NotImplementedException();
    }
}