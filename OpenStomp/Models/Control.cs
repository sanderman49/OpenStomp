namespace OpenStomp.Models;

public class Control(string name)
{
    public string Name { get; set; } = name;
    
    public bool Default { get; set; }
    
    public enum ControlType
    {
        Disabled = 0,
        Toggle = 1,
        Momentary = 2,
        Selection = 3,
        BpmTap = 4,
    }
}