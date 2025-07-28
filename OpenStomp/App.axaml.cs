using System.Reflection;
using Avalonia;
using Avalonia.Controls.ApplicationLifetimes;
using Avalonia.Markup.Xaml;
using OpenStomp.Models.Directory;
using OpenStomp.ViewModels;
using OpenStomp.Views;
using ReactiveUI;
using Splat;

namespace OpenStomp;

public partial class App : Application
{
    public override void Initialize()
    {
        AvaloniaXamlLoader.Load(this);
    }

    public override void OnFrameworkInitializationCompleted()
    {
        if (ApplicationLifetime is IClassicDesktopStyleApplicationLifetime desktop)
        {
            Locator.CurrentMutable.RegisterViewsForViewModels(Assembly.GetExecutingAssembly());
            
            desktop.MainWindow = new MainWindow
            {
                DataContext = new ShellViewModel(),
            };
            
            Config.GenerateConfigPath();
        }

        base.OnFrameworkInitializationCompleted();
    }
}