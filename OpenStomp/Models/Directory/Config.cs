using System;
using System.Collections.ObjectModel;
using System.IO;
using System.Runtime.InteropServices;
using System.Text.Json;
using Avalonia.Collections;


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
    
    public static void SaveScenes(AvaloniaList<Scene> scenes)
    {
        JsonSerializerOptions options = new()
        {
            WriteIndented = true
        };
        
        var serializedScenes = JsonSerializer.Serialize(scenes, options);

        string scenesPath = GetScenesPath();
        
        File.WriteAllText(scenesPath, serializedScenes);
    }
    
    
    public static AvaloniaList<Scene>? GetScenes()
    {
        string scenesPath = GetScenesPath();

        string serializedPrograms;
        
        try
        {
            serializedPrograms = File.ReadAllText(scenesPath);
        }
        catch (FileNotFoundException)
        {
            return null;
        }

        var scenes = JsonSerializer.Deserialize<AvaloniaList<Scene>>(serializedPrograms);

        return scenes;
    }
    
    public static string GetScenesPath()
    {
        var configPath = GetConfigPath();
        var programsPath = Path.Join(configPath, "scenes.json");

        return programsPath;
    }

    public static string GetSettingsPath()
    {
        var configPath = GetConfigPath();
        var programsPath = Path.Join(configPath, "settings.json");

        return programsPath;
    }
}