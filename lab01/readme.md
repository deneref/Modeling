﻿# Лабораторная работа 1
**Тема**: Программная реализация приближенного аналитического метода и численных
алгоритмов первого и второго порядков точности при решении задачи Коши для ОДУ.

**Цель работы.** Получение навыков решения задачи Коши для ОДУ методами Пикара и
явными методами первого порядка точности (Эйлера) и второго порядка точности (РунгеКутта).
Исходные данные.
1. ОДУ, не имеющее аналитического решения

y'(x) = x^2 + y(x)^2

y(x0) = 0

**Результат работы программы.**
1. Таблица, содержащая значения аргумента с заданным шагом в интервале [0, xmax] и
результаты расчета функции
u(x)
в приближениях Пикара (от 1-го до 4-го), а также
численными методами. Границу интервала xmax выбирать максимально возможной из
условия, чтобы численные методы обеспечивали точность вычисления решения уравнения
u(x)
до второго знака после запятой.


