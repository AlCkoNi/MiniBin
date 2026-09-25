# Release checklist

Перед публикацией новой версии:

- обновить `src/version.go`;
- обновить `VERSION.txt`;
- обновить `CHANGELOG.md`;
- обновить README и документацию, если менялось поведение;
- собрать релиз через `scripts/build_release.ps1`;
- проверить SHA-256;
- запустить сборку на Windows;
- проверить меню, иконку, Empty, Auto Clean, About и Exit;
- создать git tag `vX.Y`;
- создать GitHub Release и приложить ZIP.
