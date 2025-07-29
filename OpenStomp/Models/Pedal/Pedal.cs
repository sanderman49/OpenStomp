using System.Collections.Generic;

namespace OpenStomp.Models.Pedal;

// Represents the physical pedal.
public class Pedal
{
    public List<Button> Buttons { get; } = new();

    private Scene? _currentScene;
    public Scene? CurrentScene { get => _currentScene; }

    private int _pageNumber;
    public int PageNumber { get => _pageNumber; }

    private Display _display;
    
    public Pedal()
    {
        // Init values.
        // Hardcode in the parameters for my development pedal. Will be done dynamically later.
        for (int i = 0; i < 4; i++) // 4 Program.
        {
            Buttons.Add(new Button(Button.ButtonMode.Scene, new Led()));
        }
        for (int i = 0; i < 4; i++) // 4 Control.
        {
            Buttons.Add(new Button(Button.ButtonMode.Control, new Led()));
        }
        for (int i = 0; i < 2; i++) // 2 Nav. Don't have LEDs.
        {
            Buttons.Add(new Button(Button.ButtonMode.Navigation));
        }

        _pageNumber = 0;
        _display = new Display();
    }

    public List<Button> GetButtonsOfMode(Button.ButtonMode mode)
    {
        List<Button> buttons = new();
        
        foreach (var button in Buttons)
        {
            if (button.Mode == mode)
            {
                buttons.Add(button);
            }
        }

        return buttons;
    }
    
    public void SetCurrentScene(Scene scene)
    {
        _currentScene = scene;
        
        _display.DisplayMessage(scene.Name);
    }
}