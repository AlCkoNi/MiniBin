# Разработка и сборка MiniBin

## Требования

- Windows 10/11 для запуска и runtime-тестов.
- Go 1.23 или новее для сборки.
- Git — рекомендуется, но не требуется.

Проект использует только стандартную библиотеку Go.

## Сборка на Windows

Из корня проекта:

```powershell
.\scripts\build_release.ps1
```

Скрипт:

1. компилирует `windows/386` GUI-приложение;
2. копирует `minibin.ini` и пять ICO-файлов;
3. создаёт `dist\MiniBin - 1.1`;
4. считает SHA-256;
5. создаёт `MiniBin - 1.1.zip`.

Для простой отладочной сборки:

```powershell
$env:GOOS = "windows"
$env:GOARCH = "386"
$env:CGO_ENABLED = "0"
go build -buildvcs=false -trimpath -o build\MiniBin-debug.exe .\src
```

## Воспроизводимая release-сборка

Релиз собирается со следующими параметрами:

```text
GOOS=windows
GOARCH=386
CGO_ENABLED=0
-trimpath
-buildvcs=false
-ldflags="-s -w -H=windowsgui -buildid="
```

Это убирает локальные пути, debug symbols и Go build ID. При одинаковой версии Go и одинаковом исходном дереве результат должен быть детерминированным.

## Изменение версии

Измените:

```go
const AppVersion = "1.1"
```

в `src/version.go`, затем обновите `CHANGELOG.md`, README и документацию. Для этого проекта принята последовательная схема `1.1 -> 1.2 -> 1.3`, если отдельно не согласована другая схема.

## Проверка перед релизом

На реальной Windows проверить:

- запуск без консольного окна;
- появление иконки в трее;
- `Open`;
- disabled/enabled состояние `Empty`;
- ручную очистку;
- двойной клик;
- каждый вариант Auto Clean;
- сохранение `minibin.ini`;
- обновление пяти состояний иконки;
- `About`;
- `Exit`;
- поведение при заблокированных файлах в Temp.

После этого опубликовать ZIP как GitHub Release и приложить SHA-256.
