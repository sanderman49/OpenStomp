using System;
using System.Collections.ObjectModel;
using System.IO;
using System.Runtime.InteropServices;
using System.Text.Json;


namespace OpenStomp.Models.Directory;

public static class Config
{
    // Generate the config directory if it doesn't exist.
    public static void GenerateConfigPath()
    {
        string configPath = GetConfigPath();
        
        if (!System.IO.Directory.Exists(configPath))
        {
            System.IO.Directory.CreateDirectory(configPath);
        }
    }
    
    // Get the config directory for each OS platform.
    public static string GetConfigPath()
    {
        if (RuntimeInformation.IsOSPlatform(OSPlatform.Linux))
        {
            return Path.Join(Environment.GetFolderPath(Environment.SpecialFolder.UserProfile), ".config", "openstomp");
        }
        if (RuntimeInformation.IsOSPlatform(OSPlatform.OSX))
        {
            return Path.Join(Environment.GetFolderPath(Environment.SpecialFolder.UserProfile), "Library",
                "Application Support", "openstomp");
        }
        if (RuntimeInformation.IsOSPlatform(OSPlatform.Windows))
        {
            return Path.Join(Environment.GetFolderPath(Environment.SpecialFolder.UserProfile), "AppData", "Local", "openstomp");
        }
        if (OperatingSystem.IsAndroid())
        {
            return Path.Join(Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData));
        }

        return "";
    }

    public static string GetSettingsPath()
    {
        var configPath = GetConfigPath();
        var programsPath = Path.Join(configPath, "settings.json");

        return programsPath;
    }
}