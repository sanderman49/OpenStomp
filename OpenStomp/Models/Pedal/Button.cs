namespace OpenStomp.Models.Pedal;

// Represents a physical button on the pedal.
public class Button
{
    private bool _enabled;
    public bool Enabled { get => _enabled; }

    private Led? _led;
    
    public ButtonMode Mode { get; }

    public enum ButtonMode
    {
        Navigation = 0,
        Scene = 1,
        Control = 2,
    }
    
    public Button(ButtonMode mode, Led? led = null)
    {
        Mode = mode;
        _enabled = false;

        _led = led;
    }
    
    public void SetEnableState(bool enableState)
    {
        _led?.SetState(enableState);
    }
}