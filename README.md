### translate

translate - это [API](https://tech.yandex.ru/translate/) клиент и
cli утилита для [Яндекс.Переводчика](https://translate.yandex.ru/).

```
go get -u github.com/koorgoo/translate
```

Для работы `translate` необходимо установить значение переменной среды
`YANDEXTRANSLATEAPIKEY`, используя
[полученный](https://tech.yandex.ru/translate/doc/dg/concepts/api-keys-docpage/)
ключ.

```
translate -ls          # вывести доступные направления перевода
translate -lang hello  # определить язык текста ("en")
translate hello world  # показать перевод ("привет мир")
```


--

P.S. Спасибо Rob Pike за [вдохновение](https://github.com/robpike/translate).
