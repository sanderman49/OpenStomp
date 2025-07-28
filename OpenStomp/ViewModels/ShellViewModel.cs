using ReactiveUI;

namespace OpenStomp.ViewModels;

public class ShellViewModel : ViewModelBase, IScreen
{
    public RoutingState Router { get; }

    public ShellViewModel()
    {
        Router = new RoutingState();
        
        Router.NavigateAndReset.Execute(new MainViewModel(this));
    }
}