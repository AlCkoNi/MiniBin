# Как добавлять свои функции в MiniBin

MiniBin 1.1 специально оставлен как универсальная базовая версия. Новые функции рекомендуется делать конфигурируемыми и не привязывать к имени пользователя, букве диска или конкретной папке разработчика.

## Рекомендуемая точка расширения: My Function

Начиная с MiniBin 1.1 в главном меню постоянно присутствует пустое подменю `My Function`. Новые пользовательские функции рекомендуется добавлять именно туда:

```go
appendMenu(appState.myMenu, mfString, uintptr(idMyAction), "My Action")
```

Это сохраняет верхний уровень меню стабильным. Полная инструкция и примеры находятся в [`MY_FUNCTION_RU.md`](MY_FUNCTION_RU.md).

## 1. Добавление простого пункта меню

В `src/win32.go` добавьте новый ID:

```go
const (
    // ...
    idMyAction
)
```

В `buildMenus()` из `src/main.go` добавьте пункт:

```go
appendMenu(appState.myMenu, mfString, uintptr(idMyAction), "My Action")
```

В `handleCommand()` добавьте обработчик:

```go
case idMyAction:
    runMyAction()
```

Реализацию лучше вынести в отдельный файл, например `src/feature_myaction.go`.

## 2. Подменю

Создайте отдельный popup menu:

```go
sub := createPopupMenu()
appendMenu(sub, mfString, uintptr(idMyAction), "Run")
appendMenu(appState.myMenu, mfPopup, sub, "Tools")
```

Если подменю живёт до завершения приложения, сохраните его handle в `appState`.

## 3. Переключатель с галочкой

Добавьте поле в `Config`, загрузку и сохранение в `src/config.go`. Затем используйте:

```go
checkMenuItem(menu, idMyAction, appState.config.MyActionEnabled)
```

При выборе команды переключайте значение и сохраняйте конфигурацию.

## 4. Долгая операция

Не выполняйте тяжёлую файловую работу непосредственно внутри `windowProc`. Запускайте её в goroutine и отправляйте сообщение скрытому окну после завершения. Это сохраняет отзывчивость системного трея.

## 5. Функция, запускаемая вместе с очисткой

Если функция должна выполняться при ручной и автоматической очистке, добавьте вызов внутрь goroutine в `startCleanup()`:

```go
go func() {
    emptyRecycleBin(appState.hwnd)
    clearTemporaryFiles()
    if appState.config.MyActionEnabled {
        runMyAction()
    }
    postMessage(appState.hwnd, wmCleanDone)
}()
```

Для потенциально опасной функции безопасное значение по умолчанию должно быть `false`.

## 6. Пути

Не задавайте абсолютный путь, зависящий от конкретного компьютера. Используйте `%USERPROFILE%`, `%LOCALAPPDATA%`, `os.UserHomeDir()` либо пользовательскую настройку.

## 7. Удаление файлов

Перед добавлением рекурсивного удаления:

- ограничьте разрешённую область;
- не следуйте неожиданно за reparse point/junction;
- не удаляйте корень диска или профиль пользователя целиком;
- документируйте, что будет удалено;
- делайте опасную функцию выключенной по умолчанию.

## 8. Выпуск новой версии

1. Измените `AppVersion`.
2. Обновите `CHANGELOG.md`.
3. Обновите документацию новой функции.
4. Запустите `scripts/build_release.ps1`.
5. Проведите runtime-тест на Windows.
6. Создайте Git tag `vX.Y` и GitHub Release.
