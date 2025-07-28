using Avalonia;
using Avalonia.Controls;
using Avalonia.Controls.Primitives;
using Avalonia.ReactiveUI;
using OpenStomp.ViewModels;

namespace OpenStomp.Views;

public partial class MainView : ReactiveUserControl<MainViewModel>
{

    public MainView()
    {
        InitializeComponent();
    }
}