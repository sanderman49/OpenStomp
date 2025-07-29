using System;
using System.Collections.Generic;
using System.ComponentModel;
using Avalonia.Collections;
using OpenStomp.Models.Pedal;
using ReactiveUI;

namespace OpenStomp.Models;

public class Scene : ReactiveObject
{
    public string Id { get; set; } = Guid.NewGuid().ToString();
    public string Name { get; set; } = null!;

    public AvaloniaList<Control> Controls { get; set; } = new();

    public Scene(string name, int controlNum)
    {
        Name = name;
        
        for (int i = 0; i < controlNum; i++)
        {
            Controls.Add(new Control($"Control {i + 1}"));
        }
    }

    public Scene()
    {
        
    }
}