// auto-photo-rename.cpp : This file contains the 'main' function. Program execution begins and ends there.
//

#include <iostream>
#include <sstream>
#include <string>
#include <array>
#include <vector>
#include <algorithm>
#include <filesystem>
#include <exiv2/exiv2.hpp>

namespace limit
{
    static const uint32_t maxCmdParam{ 10 };
}

namespace cmd
{
    static const std::string exit{ "exit" };
    static const std::string quit{ "quit" };
    static const std::string help{ "help" };
    static const std::string name{ "name" };
}

static const std::vector<std::string> imgTypes{ ".jpg", ".jpeg", ".png", ".heic"};

void handleCmdExit();
void handleCmdHelp();
void handleCmdRename(std::string dirPath);

int main(int argc, char* argv[])
{
    std::cout << "1. input help to get command instruction;\n";
    std::cout << "2. input exit or quit to stop;\n";
    std::cout << "3. input name [folder path] to rename all photos in the specific folder;\n";
    std::cout << std::endl;

    std::array<std::string, limit::maxCmdParam> cmds;
    cmds.fill("");

    while (std::cin.good())
    {
        std::string input;
        std::getline(std::cin, input);

        std::istringstream iss{ input };
        uint32_t cmdCount{ 0 };
        std::string cmd;

        while (std::getline(iss, cmd, ' '))
        {
            if (!cmd.empty() && cmdCount < cmds.size())
            {
                cmds[cmdCount] = cmd;
                cmdCount++;
            }

            if (cmdCount >= limit::maxCmdParam)
            {
                break;
            }
        }

        if (cmdCount <= cmds.size())
        {
            if (cmdCount == 1)
            {
                if (cmds[0] == cmd::exit || cmds[0] == cmd::quit)
                {
                    handleCmdExit();
                }
                else if (cmds[0] == cmd::help)
                {
                    handleCmdHelp();
                }
            }
            else if (cmdCount == 2)
            {
                if (cmds[0] == cmd::name)
                {
                    handleCmdRename(cmds[1]);
                }
            }
        }

        cmds.fill("");
    }

    return 0;
}

void handleCmdExit()
{
    std::cout << "bye..." << std::endl;
    exit(0);
}

void handleCmdHelp()
{
    std::cout << "1. input help to get command instruction;\n";
    std::cout << "2. input exit or quit to stop;\n";
    std::cout << "3. input name [folder path] to rename all photos in the specific folder;\n";
    std::cout << std::endl;
}

void handleCmdRename(std::string dirPath)
{
    namespace fs = std::filesystem;

    const fs::path dir{ dirPath };
    if (!fs::exists(dir))
    {
        std::cout << "folder not exists: " << dirPath << std::endl;
        return;
    }

    if (!fs::is_directory(dir))
    {
        std::cout << "this is not a folder: " << dirPath << std::endl;
        return;
    }

    std::vector<std::string> images;

    auto transHandle = [](char c) -> char
    {
        if (std::isalpha(c) && std::isupper(c))
        {
            return std::tolower(c);
        }
        else
        {
            return c;
        }
    };

    for (const auto& entry : fs::directory_iterator(dir))
    {
        if (entry.path().has_extension())
        {
            auto extension = entry.path().extension().string();
            std::cout << extension << std::endl;
            std::transform(extension.begin(), extension.end(), extension.begin(), transHandle);
            std::cout << extension << std::endl;

            /*auto result = std::find(imgTypes.cbegin(), imgTypes.cend(), extension);

            if (result != imgTypes.cend())
            {
                images.emplace_back(entry.path().string());
            }*/
        }
    }

    for (auto& entry : images)
    {
        //std::cout << entry << std::endl;
        //try
        //{
        //    Exiv2::Image::UniquePtr image = Exiv2::ImageFactory::open(entry);
        //    if (!image.get()) {
        //        continue;
        //    }

        //    image->readMetadata();

        //    Exiv2::ExifData& exifData = image->exifData();
        //    if (exifData.empty())
        //    {
        //        std::cerr << "no EXIF data" << std::endl;
        //    }

        //    //Exiv2::Exifdatum time = exifData["Exif.Photo.DateTimeOriginal"];
        //    //if (time.count())
        //    //{
        //    //    std::cout << time.toString() << std::endl;
        //    //}

        //    //// 打印所有EXIF键值对
        //    //for (auto const& tag : exifData) {
        //    //    std::cout << tag.key() << ": " << tag.toString() << std::endl;
        //    //}
        //}
        //catch (Exiv2::Error& e)
        //{
        //    std::cerr << "Caught Exiv2 exception: " << e.what() << std::endl;
        //    return;
        //}
        //catch (...)
        //{
        //    std::cerr << "Caught unknown exception" << std::endl;
        //    return;
        //}
    }
}
