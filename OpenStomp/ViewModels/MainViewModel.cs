using System;
using System.Collections.Generic;
using System.Linq;
using System.Windows.Input;
using Avalonia;
using Avalonia.Collections;
using Avalonia.Controls.Platform;
using OpenStomp.Models;
using OpenStomp.Models.Directory;
using OpenStomp.Models.Pedal;
using ReactiveUI;

namespace OpenStomp.ViewModels;

public class MainViewModel : ViewModelBase, IRoutableViewModel
{
    public IScreen HostScreen { get; }

    public string? UrlPathSegment { get; } = Guid.NewGuid().ToString().Substring(0, 5);
    
    public Pedal Pedal { get; set; }

    private int _page;
    public int Page
    {
        get => _page;
        set
        {
            this.RaiseAndSetIfChanged(ref _page, value);
            this.RaisePropertyChanged(nameof(FormattedPage));
        }
    }

    public string FormattedPage => (Page + 1).ToString();

    private int _sceneLimit;
    private int _controlsPerScene;
    
    private int _pageLimit;
    
    public AvaloniaList<Scene> Scenes { get; set; }
    
    private int _scenesPerPage;

    private AvaloniaList<Scene> _visibleScenes;
    public AvaloniaList<Scene> VisibleScenes
    {
        get => _visibleScenes;
        set => this.RaiseAndSetIfChanged(ref _visibleScenes, value);
    }


    private AvaloniaList<Control> _visibleControls;
    public AvaloniaList<Control> VisibleControls
    {
        get => _visibleControls;
        set => this.RaiseAndSetIfChanged(ref _visibleControls, value);
    }

    private Scene _selectedScene;
    public Scene SelectedScene
    {
        get => _selectedScene;
        set
        {
            this.RaiseAndSetIfChanged(ref _selectedScene, value);
        }
    }

    public ICommand SelectScene { get; }
    public ICommand EditControl { get; }
    public ICommand IncrementPageNumber { get; }
    public ICommand DecrementPageNumber { get; }
    
    public MainViewModel(IScreen screen)
    {
        // Init initial values
        HostScreen = screen;
        
        Pedal = new Pedal();
        Page = 0;
        _sceneLimit = 48; // Arbitrary max scene number.
        
        List<Button> sceneButtons = Pedal.GetButtonsOfMode(Button.ButtonMode.Scene); // All physical scene buttons.
        List<Button> controlButtons = Pedal.GetButtonsOfMode(Button.ButtonMode.Control); // All physical control buttons.
        
        _scenesPerPage = sceneButtons.Count; // The number of physical scene buttons.
        _controlsPerScene = controlButtons.Count; // The number of physical control buttons.
        _pageLimit = _sceneLimit / _scenesPerPage; // Number of pages to need to fit all scenes.

        // Init scenes
        var scenes = Config.GetScenes();

        if (scenes == null)
            GenerateAndSaveScenes();
        else
            Scenes = scenes;
        
        SelectedScene = Scenes[0];
        VisibleControls = SelectedScene.Controls;
        
        VisibleScenes = new();
        SetVisibleScenesForPage(0);
        
        // Init ICommands
        SelectScene = ReactiveCommand.Create<Scene>((scene) =>
        {
            SelectedScene = scene;
            VisibleControls = scene.Controls;
        });
        
        EditControl = ReactiveCommand.Create<Control>((control) =>
        {
            Console.WriteLine(control.Name);
        });
        
        IncrementPageNumber = ReactiveCommand.Create(() => ModifyPageNumber(1));
        DecrementPageNumber = ReactiveCommand.Create(() => ModifyPageNumber(-1));
    }

    private void ModifyPageNumber(int number)
    {
        // Don't let page be less than 0 or greater than the page limit.
        if ((Page + number) >= 0 && (Page + number) < _pageLimit)
            Page += number;

        VisibleScenes.Clear();
        
        SetVisibleScenesForPage(Page);
    }

    private void SetVisibleScenesForPage(int page)
    {
        for (int i = _scenesPerPage * page; i < (_scenesPerPage * page) + _scenesPerPage; i++)
        {
            VisibleScenes.Add(Scenes[i]);
        }
    }

    private void GenerateAndSaveScenes()
    {
        Scenes = new();
        
        for (int i = 0; i < _sceneLimit; i++)
        {
            Scenes.Add(new Scene($"Scene {i + 1}", _controlsPerScene));
        }
        
        Config.SaveScenes(Scenes);
    }
    
}