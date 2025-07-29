using System;
using ReactiveUI;

namespace OpenStomp.Models;

public class Control(string name) : ReactiveObject
{
    public string Id { get; set; } = Guid.NewGuid().ToString();
    
    private string _name = name;
    public string Name
    {
        get => _name;
        set => this.RaiseAndSetIfChanged(ref _name, value);
    }

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