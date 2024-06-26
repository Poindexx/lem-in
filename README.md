# lem-in - Digital Ant Farm Simulation

## Цели

Этот проект представляет собой программу lem-in, которая считывает данные из файла (описывающего муравейник и муравьев) в аргументах.

После успешного нахождения самого быстрого пути, lem-in отображает содержимое файла, а также каждое движение муравьев из комнаты в комнату.

## Как это работает?

- Вы создаете муравейник с туннелями и комнатами.
- Вы помещаете муравьев с одной стороны и наблюдаете, как они находят выход.
- Цель - довести всех муравьев из комнаты ##start до комнаты ##end за наименьшее количество ходов.
- Самый короткий путь не всегда является самым простым.
- Некоторые муравейники могут иметь много комнат и туннелей, но без пути между ##start и ##end.
- В некоторых муравейниках комнаты могут ссылаться сами на себя, что делает поиск пути бесконечным. Также могут быть другие недопустимые или плохо отформатированные входные данные.
- В таких случаях программа должна вернуть сообщение об ошибке "ERROR: invalid data format". Можно уточнить сообщение об ошибке, например: "ERROR: invalid data format, invalid number of Ants" или "ERROR: invalid data format, no start room found".
- Результаты должны отображаться в стандартном выводе в следующем формате:
```sh
number_of_ants
the_rooms
the_links

Lx-y Lz-w Lr-o ...
```

## Использование

Примеры использования:

```sh
$ go run . test0.txt
3
##start
1 23 3
2 16 7
3 16 3
4 16 5
5 9 3
6 1 5
7 4 8
##end
0 9 5
0-4
0-6
1-3
4-3
5-2
3-5
4-2
2-1
7-6
7-2
7-4
6-5

L1-3 L2-2
L1-4 L2-5 L3-3
L1-0 L2-6 L3-4
L2-0 L3-0
```

##Используемые пакеты
Только стандартные пакеты Go.

Инструкции
Создавайте туннели и комнаты.
Комната никогда не начинается с буквы L или # и не должна содержать пробелов.
Соединяйте комнаты между собой с любым количеством туннелей.
Каждый туннель соединяет только две комнаты, никогда больше.
Комнату можно соединить с несколькими другими комнатами.
Две комнаты не могут иметь более одного туннеля, соединяющего их.
Каждая комната может содержать только одного муравья одновременно (кроме ##start и ##end, где может быть любое количество муравьев).
Каждый туннель может использоваться только один раз за ход.
Чтобы быть первым, кто прибывает, муравьи должны выбирать самый короткий путь или пути. Они также должны избегать пробок и ходить по своим собратьям.
Выводить только муравьев, которые сделали ход на каждом ходу, и каждый муравей может двигаться только один раз и через туннель (комната на принимающем конце должна быть пустой).
Имена комнат не обязательно будут числами и в порядке.
Любая неизвестная команда должна игнорироваться.
Программа должна тщательно обрабатывать ошибки и ни в коем случае не завершаться непредвиденным образом.
Координаты комнат всегда будут целыми числами.
