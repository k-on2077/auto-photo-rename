#include "FuncLib.h"

void func::transform2Lower(std::string& target)
{
    auto handle = [](char c) -> char { return (std::isalpha(c) && std::isupper(c)) ? std::tolower(c) : c; };
    std::transform(target.begin(), target.end(), target.begin(), handle);
}
