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
'test2main' -> overwrite main with test toml
'test2local' -> overwrite local with test toml

'read' followed by:
    'test' -> print test toml contents
    'local' -> print local dotstrike toml contents
    'main' or 'global' -> print main toml contents
'wipe' followed by:
    'test' -> clear test file data
"""
listnames = ['test', 'local', 'global', 'main']


def wipe_file(fp: Path) -> bool:
    if fp.exists():
        with fp.open("r+") as file:
            _ = file.seek(0)
            size = file.truncate(0)
            file.close()
        if size == 0:
            return True
    return False


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


def try_copy_twice(src: Path, src_abs: Path, dest: Path):
    if src.exists():
        copy_file(src, dest)
    elif src_abs.exists():
        copy_file(src_abs, dest)
    else:
        print(f'no file at {src} or {src_abs}')


def linenum_print(txt: str):
    tln: list[str] = txt.split("\n")
    for i in range(len(tln)):
        tln[i] = f'{i+1}|' + tln[i]
        print(tln[i])


def main():
    if len(sys.argv) <= 1:
        print(helpstring)
    else:
        a = sys.argv[1].lower()
        if a == 'read':
            if len(sys.argv) > 2 and sys.argv[2].lower() in listnames:  # noqa
                match sys.argv[2].lower():
                    case 'test':
                        linenum_print(readf(testtoml))
                    case 'local':
                        linenum_print(readf(localtoml))
                    case 'main' | 'global':
                        linenum_print(readf(maintoml))
                    case _:
                        print('unknown file')
        elif a == 'wipe':
            if len(sys.argv) > 2 and sys.argv[2].lower() in listnames:
                match sys.argv[2].lower():
                    case 'test':
                        if testtoml.exists() and wipe_file(testtoml):
                            print("success")
                        else:
                            print("outcome unknown: data may still exist")
                    case 'local':
                        print("not set up to wipe local")
                    case 'main' | 'global':
                        print("not set up to wipe main")
                    case _:
                        print('unknown file')
        else:
            match a:
                case 'cleartest' | 'test' | 'overwritetest':
                    try_copy_twice(localtoml, localtoml_abs, testtoml)
                case 'push' | 'to_main':
                    try_copy_twice(localtoml, localtoml_abs,  maintoml)
                case 'get' | 'pull':
                    if maintoml.exists():
                        copy_file(maintoml, localtoml)
                    else:
                        print(f'no file at {maintoml}')
                case 'test-to-main' | 'test2main' | 'test to main':
                    try_copy_twice(testtoml, testtoml_abs, maintoml)
                case 'test-to-local' | 'test2local' | 'test to local':
                    try_copy_twice(testtoml, testtoml_abs, localtoml)
                case _:
                    print('unknown arg')
                    print(helpstring)


if __name__ == "__main__":
    main()
# %%
