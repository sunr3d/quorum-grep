#!/bin/bash

echo "Сравнение результатов grep и mygrep"

echo "==Тест 1: Стандартный поиск=="
echo "=GREP=:"
grep "pattern" test.txt
echo "=MYGREP=:"
../mygrep "pattern" test.txt

echo "==Тест 2: Поиск с флагом -n (номера строк)=="
echo "=GREP=:"
grep -n "pattern" test.txt
echo "=MYGREP=:"
../mygrep -n "pattern" test.txt

echo "==Тест 3: Поиск с флагом -A (строки после)=="
echo "=GREP=:"
grep -A 2 "pattern" test.txt
echo "=MYGREP=:"
../mygrep -A 2 "pattern" test.txt

echo "==Тест 4: Поиск с флагом -B (строки перед)=="
echo "=GREP=:"
grep -B 2 "pattern" test.txt
echo "=MYGREP=:"
../mygrep -B 2 "pattern" test.txt

echo "==Тест 5: Поиск с флагом -C (строки вокруг)=="
echo "=GREP=:"
grep -C 2 "pattern" test.txt
echo "=MYGREP=:"
../mygrep -C 2 "pattern" test.txt

echo "==Тест 6: Поиск с флагом -c (количество строк)=="
echo "=GREP=:"
grep -c "pattern" test.txt
echo "=MYGREP=:"
../mygrep -c "pattern" test.txt

echo "==Тест 7: Поиск с флагом -i (игнорировать регистр)=="
echo "=GREP=:"
grep -i "PATTERN" test.txt
echo "=MYGREP=:"
../mygrep -i "PATTERN" test.txt

echo "==Тест 8: Поиск с флагом -v (инвертировать результат)=="
echo "=GREP=:"
grep -v "pattern" test.txt
echo "=MYGREP=:"
../mygrep -v "pattern" test.txt

echo "==Тест 9: Поиск с флагом -F (фиксированная строка)=="
echo "=GREP=:"
grep -F "special.pattern" test.txt
echo "=MYGREP=:"
../mygrep -F "special.pattern" test.txt

echo "==Тест 10: Большой файл (1000+ строк)=="
echo "=GREP=:"
grep "pattern" big_test.txt | head -5
echo "=MYGREP=:"
../mygrep "pattern" big_test.txt | head -5

echo "==Тест 11: Контекст на большом файле -A 2 -B 1=="
echo "=GREP=:"
grep -A 2 -B 1 "pattern" big_test.txt | head -10
echo "=MYGREP=:"
../mygrep -A 2 -B 1 "pattern" big_test.txt | head -10

echo "==Тест 12: Поиск в stdin (пайп)=="
echo "=GREP=:"
echo -e "line1\npattern found\nline3\nanother pattern\nline5" | grep "pattern"
echo "=MYGREP=:"
echo -e "line1\npattern found\nline3\nanother pattern\nline5" | ../mygrep "pattern"

echo "==Тест 13: Поиск в нескольких файлах=="
echo "=GREP=:"
grep "pattern" test.txt context_test.txt
echo "=MYGREP=:"
../mygrep "pattern" test.txt context_test.txt

echo "==Тест 14: Флаг -C (around) на большом файле=="
echo "=GREP=:"
grep -C 1 "pattern" big_test.txt | head -10
echo "=MYGREP=:"
../mygrep -C 1 "pattern" big_test.txt | head -10

echo "==Тест 15: Производительность (время выполнения)=="
echo "=GREP=:"
time grep -c "pattern" big_test.txt
echo "=MYGREP=:"
time ../mygrep -c "pattern" big_test.txt

echo "==Тест 16: Регулярные выражения=="
echo "=GREP=:"
grep "pattern[0-9]" big_test.txt | head -5
echo "=MYGREP=:"
../mygrep "pattern[0-9]" big_test.txt | head -5

echo "==Тест 17: Фиксированная строка на большом файле=="
echo "=GREP=:"
grep -F "pattern found here" big_test.txt
echo "=MYGREP=:"
../mygrep -F "pattern found here" big_test.txt

echo "Конец тестов..."