# %%
from os import getcwd
from pathlib import Path
from shutil import copyfile
import sys

localtoml_abs = Path(r'D:\coding\github\dotstrike\[samplefiles]\dotstrikeData.toml')  # noqa
testtoml_abs = Path(r'D:\coding\github\dotstrike\[samplefiles]\test_dotstrikeData.toml')  # noqa

localtoml = Path(getcwd()).joinpath('[samplefiles]\\dotstrikeData.toml')
maintoml: Path = Path.home().joinpath(r'.config\dotstrike\dotstrikeData.toml')
testtoml = Path(getcwd()).joinpath('[samplefiles]\\test_dotstrikeData.toml')

helpstring = """----[need arg:]----
'test' or 'cleartest' -> overwrite test file
'push' -> dotstrike local toml pushed to main toml ('~\\.config\\dotstrike\\')
'get' or 'pull' -> overwrite dotstrike local toml with main toml

'read' followed by:
    'test' -> print test toml contents
    'local' -> print local dotstrike toml contents
    'main' or 'global' -> print main toml contents
"""


def copy_file(source: Path | str, dest: Path | str) -> None:
    sp: Path = Path(source)
    dp: Path = Path(dest)

    if sp and dp and sp.exists() and sp.is_file():
        outpath: Path = copyfile(src=sp.absolute(), dst=dp.absolute())
        if outpath == dp or outpath == dest:
            print(f"successful: {outpath} overwritten")
        else:
            print(f"unknown outcome: outpath = {outpath}")


def readf(f: Path) -> str:
    if f.exists() and f.is_file():
        with f.open() as rf:
            return str(rf.read())
    return ""


def main():
    if len(sys.argv) <= 1:
        print(helpstring)
    a = sys.argv[1].lower()
    if a == 'read':
        if len(sys.argv) > 2 and sys.argv[2] in ['test', 'local', 'main', 'global']:  # noqa
            match sys.argv[2].lower():
                case 'test':
                    print(readf(testtoml))
                case 'local':
                    print(readf(localtoml))
                case 'main' | 'global':
                    print(readf(maintoml))
                case _:
                    print('unknown file')
    else:
        match a:
            case 'cleartest' | 'test' | 'overwritetest':
                if localtoml.exists():
                    copy_file(localtoml, testtoml)
                elif localtoml_abs.exists():
                    copy_file(localtoml_abs, testtoml)
                else:
                    print(f'no file at {localtoml} or {localtoml_abs}')
            case 'push' | 'to_main':
                if localtoml.exists():
                    copy_file(localtoml, maintoml)
                elif localtoml_abs.exists():
                    copy_file(localtoml_abs, testtoml)
                else:
                    print(f'no file at {localtoml} or {localtoml_abs}')
            case 'get' | 'pull':
                if maintoml.exists():
                    copy_file(localtoml, maintoml)
                else:
                    print(f'no file at {maintoml}')
            case _:
                print('unknown arg')
                print(helpstring)


if __name__ == "__main__":
    main()
# %%
