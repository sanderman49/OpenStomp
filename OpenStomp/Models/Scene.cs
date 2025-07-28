using System.Collections.Generic;
using OpenStomp.Models.Pedal;

namespace OpenStomp.Models;

public class Scene
{
    public string Name { get; set; } = null!;

    public List<Control> Controls { get; set; } = new();

    public Scene(string name, int controlNum)
    {
        Name = name;
        
        for (int i = 0; i < controlNum; i++)
        {
            Controls.Add(new Control($"Control {i + 1}"));
        }
    }
}