# README #

This README would normally document whatever steps are necessary to get your application up and running.




### Income pocket ###

// Клиент запрашивает визуализацию обекта по ID
     0    , 1 2 3 4  
     5    , objectID 
<-header->,<-  ID  ->

//Клиент указывает данные контроллеров
     0    ,    1     ,    2   3    
     6    ,   keys   , mouse angle 
<-header->,<-keySet->,<- radian  ->




### Output pocket ###

//Указать клиенту его ID
     0    , 1 2 3 4  
     0    , objectID 
<-header->,<-  ID  ->

//Указать текущее положение объекта
     0    , 1 2 3 4  ,  5 6 7  ,  8 9 10 ,  11 12   ,   13  14     
     1    , objectID ,  coordX ,  coordY , rotation ,tower rotation
<-header->,<-  ID  ->,<-coord->,<-coord->,<-radian->,<-  radian  ->

//Указать описание объекта
     0    , 1 2 3 4  ,    5     
     2    , objectID ,tower type
<-header->,<-  ID  ->,<- byte ->

//Указать отображение объекта ammo
     0    , 1 2 3 4 ,    5     ,  6 7 8  ,  9 10 11 ,  12 13   
     3    ,  ammoID , ammo type,  coordX ,  coordY  , rotation 
<-header->,<- ID  ->,<- byte ->,<-coord->,<-coord ->,<-radian->

//Дрон был уничтожен
     0    , 1 2 3 4  
     4    , objectID 
<-header->,<-  ID  ->

//Ammo столкнулся
     0    , 1 2 3 4  ,  5 6 7  ,  8 9 10 
     12   , objectID ,  coordX ,  coordY 
<-header->,<-  ID  ->,<-coord->,<-coord->

//Указать информацию чанка
     0    ,   1 - 1024   
     14   ,  chank data  
<-header->,<-chankRange->





### Data types ###

<-header-> (байт) указатель типа сообщения
- по заголовку определяется длина сообщения в пакете и метод его десериализации

<-  ID  -> (unsinged integer 32-bit) ID игрового объекта в кодировке BigEndian
- уникальный идентификатор игрового объекта

<-coord-> кордината объекта. Первый байт указывает chank(координата DIV 32) (0-255), второй указывает cell((координата MOD 16384) DIV 64), третий указывает точное положение внутри cell(((координата MOD 16384) MOD 64) MUL 4)
- координаты в игровом мире указаны дробным числом в пикселях

<-radian-> первый байт указывает целую часть числа (радиана DIV 1), второй указывает дробную по сотые (((радиана MOD 1) * 100) DIV 1)
- углы в игровом мире указаны дробным числом в радианах

<-chankRange-> массив данных о чанке 32х32 Cells
- данные о территории 1024 байта