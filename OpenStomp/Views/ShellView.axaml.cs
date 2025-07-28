using Avalonia;
using Avalonia.Controls;
using Avalonia.Markup.Xaml;
using Avalonia.ReactiveUI;
using OpenStomp.ViewModels;

namespace OpenStomp.Views;

public partial class ShellView : ReactiveUserControl<ShellViewModel>
{
    public ShellView()
    {
        InitializeComponent();
    }
}