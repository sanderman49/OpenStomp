using System;
using ReactiveUI;

namespace OpenStomp.ViewModels;

public class MainViewModel : ViewModelBase, IRoutableViewModel
{
    public IScreen HostScreen { get; }

    public string? UrlPathSegment { get; } = Guid.NewGuid().ToString().Substring(0, 5);

    public MainViewModel(IScreen screen)
    {
        HostScreen = screen;

    }
}