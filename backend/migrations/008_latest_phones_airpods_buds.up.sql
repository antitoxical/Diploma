-- Актуальные модели (ориентир — линейка 2025–2026 гг.): iPhone 17 / Air, AirPods Pro 3 / 4 / Max 2, Galaxy Buds4.
-- Цены в копейках BYN (100 = 1 Br). Идемпотентные INSERT.

-- Упрощённый сид иногда оставляет только корень «Наушники» (b002) без b010/b020/b021.
-- Полный 002: TWS (b020) и полноразмерные (b021) висят от «Наушники» (b010) под «Аудио» (b002).
-- Создаём b020/b021 при отсутствии, подбирая родителя: b010 если есть, иначе b002.
INSERT INTO categories (id, name, parent_id, image_url)
SELECT 'b0000000-0000-0000-0000-000000000020',
       'TWS-наушники',
       COALESCE(
           (SELECT id FROM categories WHERE id = 'b0000000-0000-0000-0000-000000000010' LIMIT 1),
           (SELECT id FROM categories WHERE id = 'b0000000-0000-0000-0000-000000000002' LIMIT 1)
       ),
       'https://placehold.co/400x300?text=TWS'
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE id = 'b0000000-0000-0000-0000-000000000020');

INSERT INTO categories (id, name, parent_id, image_url)
SELECT 'b0000000-0000-0000-0000-000000000021',
       'Полноразмерные',
       COALESCE(
           (SELECT id FROM categories WHERE id = 'b0000000-0000-0000-0000-000000000010' LIMIT 1),
           (SELECT id FROM categories WHERE id = 'b0000000-0000-0000-0000-000000000002' LIMIT 1)
       ),
       'https://placehold.co/400x300?text=OverEar'
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE id = 'b0000000-0000-0000-0000-000000000021');

INSERT INTO products (id, category_id, firm_id, name, description, price_kopeks, stock_quantity, is_active)
SELECT 'c0000000-0000-0000-0000-000000000024',
       'b0000000-0000-0000-0000-000000000001',
       'a0000000-0000-0000-0000-000000000001',
       'iPhone 17 Pro',
$p$
Флагман Apple с чипом A19 Pro, продвинутой камерой Fusion и дисплеем ProMotion.

Характеристики:
• Дисплей: Super Retina XDR, ProMotion 120 Гц, яркость для HDR, защитное стекло Ceramic Shield 2
• Платформа: Apple A19 Pro, Neural Engine для Apple Intelligence
• Камеры: система Pro (основной высокого разрешения, ультраширокий, телефото); запись видео высокого класса, режимы для слабого освещения
• Корпус: премиальные материалы серии Pro, класс влагозащиты IP68
• Связь: 5G, Wi‑Fi и Bluetooth актуальных поколений, NFC
• Память: варианты 256 ГБ и выше (уточняйте партию на складе)
• Разъём: USB‑C для зарядки и аксессуаров

Комплектация и точные характеристики модификации зависят от поставки.
$p$,
469990, 20, true
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 'c0000000-0000-0000-0000-000000000024');

INSERT INTO products (id, category_id, firm_id, name, description, price_kopeks, stock_quantity, is_active)
SELECT 'c0000000-0000-0000-0000-000000000025',
       'b0000000-0000-0000-0000-000000000001',
       'a0000000-0000-0000-0000-000000000001',
       'iPhone 17',
$p$
Основная модель линейки iPhone 17 с чипом A19 и камерой Fusion 48 Мп.

Характеристики:
• Дисплей: Super Retina XDR с ProMotion, высокая яркость для контента HDR
• Платформа: Apple A19, поддержка Apple Intelligence (функции зависят от региона и ПО)
• Камера: 48 Мп Fusion с оптическим качеством зума 2×; ультраширокий модуль
• Фронтальная камера: Center Stage и улучшенные селфи/видеозвонки
• Корпус: Ceramic Shield 2, IP68
• Память: от 256 ГБ в базовых комплектациях (по данным Apple для линейки)
• Цвета: зависят от наличия

Оптимальный выбор, если не нужен телефото «второго поколения» линейки Pro.
$p$,
339990, 28, true
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 'c0000000-0000-0000-0000-000000000025');

INSERT INTO products (id, category_id, firm_id, name, description, price_kopeks, stock_quantity, is_active)
SELECT 'c0000000-0000-0000-0000-000000000026',
       'b0000000-0000-0000-0000-000000000001',
       'a0000000-0000-0000-0000-000000000001',
       'iPhone Air',
$p$
Сверхтонкий iPhone с дисплеем ~6,5" и чипом A19 Pro — акцент на портативность и премиальный дизайн.

Характеристики:
• Толщина и масса: минимальный профиль корпуса в линейке (ориентир ~5,6 мм / ~165 г — по данным Apple)
• Дисплей: большая диагональ, высокая яркость, поддержка HDR
• Платформа: A19 Pro, графика для игр и Apple Intelligence
• Камера: 48 Мп Fusion; фронтальная 18 Мп Center Stage
• Защита: обновлённое стекло и защита от брызг/воды по классу устройства
• Разъём: USB‑C

Подходит тем, кто хочет максимально лёгкий полноразмерный iPhone без компромисса по чипу Pro.
$p$,
399990, 15, true
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 'c0000000-0000-0000-0000-000000000026');

INSERT INTO products (id, category_id, firm_id, name, description, price_kopeks, stock_quantity, is_active)
SELECT 'c0000000-0000-0000-0000-000000000027',
       'b0000000-0000-0000-0000-000000000020',
       'a0000000-0000-0000-0000-000000000001',
       'AirPods Pro 3',
$p$
Флагманские внутриканальные AirPods с усиленным ANC и датчиком пульса во время тренировок.

Характеристики:
• Чип: Apple H2 (новое поколение для серии Pro)
• Шумоподавление: до ~2× эффективнее относительно AirPods Pro 2 (по данным Apple)
• Здоровье: измерение пульса во время активности (не медицинское устройство)
• Прозрачность и адаптивное аудио; улучшенная посадка за счёт комплекта вкладышей
• Аккумулятор: до ~8 ч прослушивания с ANC (условия зависят от громкости и кодека)
• Кейс: USB‑C, MagSafe / Qi; поиск в Find My
• Защита: повышенная степень по сравнению с предыдущим поколением (уточняйте по спецификации партии)

Совместимость: лучший опыт на экосистеме Apple с последними версиями iOS / iPadOS / macOS.
$p$,
92990, 55, true
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 'c0000000-0000-0000-0000-000000000027');

INSERT INTO products (id, category_id, firm_id, name, description, price_kopeks, stock_quantity, is_active)
SELECT 'c0000000-0000-0000-0000-000000000028',
       'b0000000-0000-0000-0000-000000000020',
       'a0000000-0000-0000-0000-000000000001',
       'AirPods 4 с активным шумоподавлением',
$p$
Открытая форма «затычки» с чипом H2 и поддержкой ANC — без силиконовых насадок.

Характеристики:
• Чип: H2 с улучшенным звуком и микрофонами для звонков
• ANC: активное шумоподавление в открытой акустической конструкции
• Автономность: до ~5 ч прослушивания на цикл (без кейса; с кейсом дольше)
• Кейс: компактный, беспроводная зарядка (Qi / MagSafe по версии)
• Защита: класс пылевлагозащиты для наушников и кейса (по спецификации модели)

Удобный компромисс между ценой AirPods Pro и полным отказом от ANC.
$p$,
57990, 70, true
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 'c0000000-0000-0000-0000-000000000028');

INSERT INTO products (id, category_id, firm_id, name, description, price_kopeks, stock_quantity, is_active)
SELECT 'c0000000-0000-0000-0000-000000000029',
       'b0000000-0000-0000-0000-000000000021',
       'a0000000-0000-0000-0000-000000000001',
       'AirPods Max 2',
$p$
Накладные наушники Apple с чипом H2 и усиленным шумоподавлением.

Характеристики:
• Акустика: излучатели Apple с расширенным динамическим диапазоном
• ANC и режим прозрачности; разговоры — Voice Isolation / Adaptive Audio (по версии ПО)
• Потоковое аудио высокого разрешения по USB‑C (до 24 бит / 48 кГц — по данным Apple)
• Микрофоны: улучшенная работа в звонках и запись «студийного» качества (функции зависят от софта)
• Элементы управления: Digital Crown и кнопка управления шумом
• Автономность: до нескольких десятков часов прослушивания (без ANC / с ANC — см. спецификацию)
• Цвета: midnight, starlight, orange, purple, blue (наличие уточняйте)

Для тех, кто предпочитает полноразмерный комфорт экосистемы Apple.
$p$,
184990, 18, true
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 'c0000000-0000-0000-0000-000000000029');

INSERT INTO products (id, category_id, firm_id, name, description, price_kopeks, stock_quantity, is_active)
SELECT 'c0000000-0000-0000-0000-000000000030',
       'b0000000-0000-0000-0000-000000000020',
       'a0000000-0000-0000-0000-000000000002',
       'Samsung Galaxy Buds4 Pro',
$p$
Флагманские TWS Samsung с канальной посадкой, Hi‑Res и Galaxy AI.

Характеристики:
• Акустика: усиленный драйвер (woofer с увеличенной эффективной площадью ~на 20% к пред. поколению — Samsung)
• Hi‑Res: поддержка 24 бит / до 96 кГц на совместимых устройствах
• ANC: адаптивное шумоподавление под окружение
• Микрофоны: до шести для HD Voice и широкополосных звонков
• Galaxy AI: сценарии с переводом и голосовым управлением при связке с поддерживаемыми Galaxy (версия прошивки и регион)
• Управление: жесты головой для вызовов и Bixby (где включено)
• Кейс: компактный «ракушка», премиальная отделка; беспроводная зарядка

Лучше всего раскрываются со смартфонами Galaxy последних серий.
$p$,
77990, 42, true
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 'c0000000-0000-0000-0000-000000000030');

INSERT INTO products (id, category_id, firm_id, name, description, price_kopeks, stock_quantity, is_active)
SELECT 'c0000000-0000-0000-0000-000000000031',
       'b0000000-0000-0000-0000-000000000020',
       'a0000000-0000-0000-0000-000000000002',
       'Samsung Galaxy Buds4',
$p$
Открытая посадка Galaxy Buds4 — комфорт на весь день и Hi‑Res звук.

Характеристики:
• Конструкция: open-fit без глубокого затыкания канала — меньше утомляемости
• Звук: 24-бит Hi‑Res на совместимых телефонах Samsung
• ANC: адаптивное (уровень зависит от посадки и окружения)
• Микрофоны и HD Voice для чётких звонков
• Galaxy AI (перевод, сценарии hands-free — при поддержке телефона и региона)
• Кейс: прозрачная или полупрозрачная «ракушка» в зависимости от версии
• Цвета: White, Black, Pink Gold — уточняйте наличие

Подходят тем, кто не любит вкладыши с силиконовыми насадками, но хочет свежую серию Buds4.
$p$,
46990, 50, true
WHERE NOT EXISTS (SELECT 1 FROM products WHERE id = 'c0000000-0000-0000-0000-000000000031');

INSERT INTO product_images (product_id, image_url, is_main, sort_order)
SELECT v.id, v.url, true, 0
FROM (VALUES
    ('c0000000-0000-0000-0000-000000000024'::uuid, 'https://placehold.co/400x400?text=iPhone+17+Pro'),
    ('c0000000-0000-0000-0000-000000000025'::uuid, 'https://placehold.co/400x400?text=iPhone+17'),
    ('c0000000-0000-0000-0000-000000000026'::uuid, 'https://placehold.co/400x400?text=iPhone+Air'),
    ('c0000000-0000-0000-0000-000000000027'::uuid, 'https://placehold.co/400x400?text=AirPods+Pro+3'),
    ('c0000000-0000-0000-0000-000000000028'::uuid, 'https://placehold.co/400x400?text=AirPods+4+ANC'),
    ('c0000000-0000-0000-0000-000000000029'::uuid, 'https://placehold.co/400x400?text=AirPods+Max+2'),
    ('c0000000-0000-0000-0000-000000000030'::uuid, 'https://placehold.co/400x400?text=Galaxy+Buds4+Pro'),
    ('c0000000-0000-0000-0000-000000000031'::uuid, 'https://placehold.co/400x400?text=Galaxy+Buds4')
) AS v(id, url)
WHERE NOT EXISTS (
    SELECT 1 FROM product_images pi WHERE pi.product_id = v.id AND pi.is_main = true
);

INSERT INTO product_tags (product_id, tag_id)
SELECT 'c0000000-0000-0000-0000-000000000024', 'd0000000-0000-0000-0000-000000000001'
WHERE EXISTS (SELECT 1 FROM products WHERE id = 'c0000000-0000-0000-0000-000000000024')
  AND NOT EXISTS (SELECT 1 FROM product_tags WHERE product_id = 'c0000000-0000-0000-0000-000000000024' AND tag_id = 'd0000000-0000-0000-0000-000000000001');

INSERT INTO product_tags (product_id, tag_id)
SELECT 'c0000000-0000-0000-0000-000000000027', 'd0000000-0000-0000-0000-000000000001'
WHERE EXISTS (SELECT 1 FROM products WHERE id = 'c0000000-0000-0000-0000-000000000027')
  AND NOT EXISTS (SELECT 1 FROM product_tags WHERE product_id = 'c0000000-0000-0000-0000-000000000027' AND tag_id = 'd0000000-0000-0000-0000-000000000001');

INSERT INTO product_tags (product_id, tag_id)
SELECT 'c0000000-0000-0000-0000-000000000030', 'd0000000-0000-0000-0000-000000000001'
WHERE EXISTS (SELECT 1 FROM products WHERE id = 'c0000000-0000-0000-0000-000000000030')
  AND NOT EXISTS (SELECT 1 FROM product_tags WHERE product_id = 'c0000000-0000-0000-0000-000000000030' AND tag_id = 'd0000000-0000-0000-0000-000000000001');

INSERT INTO product_tags (product_id, tag_id)
SELECT 'c0000000-0000-0000-0000-000000000025', 'd0000000-0000-0000-0000-000000000001'
WHERE EXISTS (SELECT 1 FROM products WHERE id = 'c0000000-0000-0000-0000-000000000025')
  AND NOT EXISTS (SELECT 1 FROM product_tags WHERE product_id = 'c0000000-0000-0000-0000-000000000025' AND tag_id = 'd0000000-0000-0000-0000-000000000001');

INSERT INTO product_tags (product_id, tag_id)
SELECT 'c0000000-0000-0000-0000-000000000026', 'd0000000-0000-0000-0000-000000000001'
WHERE EXISTS (SELECT 1 FROM products WHERE id = 'c0000000-0000-0000-0000-000000000026')
  AND NOT EXISTS (SELECT 1 FROM product_tags WHERE product_id = 'c0000000-0000-0000-0000-000000000026' AND tag_id = 'd0000000-0000-0000-0000-000000000001');

INSERT INTO product_tags (product_id, tag_id)
SELECT 'c0000000-0000-0000-0000-000000000028', 'd0000000-0000-0000-0000-000000000001'
WHERE EXISTS (SELECT 1 FROM products WHERE id = 'c0000000-0000-0000-0000-000000000028')
  AND NOT EXISTS (SELECT 1 FROM product_tags WHERE product_id = 'c0000000-0000-0000-0000-000000000028' AND tag_id = 'd0000000-0000-0000-0000-000000000001');

INSERT INTO product_tags (product_id, tag_id)
SELECT 'c0000000-0000-0000-0000-000000000029', 'd0000000-0000-0000-0000-000000000001'
WHERE EXISTS (SELECT 1 FROM products WHERE id = 'c0000000-0000-0000-0000-000000000029')
  AND NOT EXISTS (SELECT 1 FROM product_tags WHERE product_id = 'c0000000-0000-0000-0000-000000000029' AND tag_id = 'd0000000-0000-0000-0000-000000000001');

INSERT INTO product_tags (product_id, tag_id)
SELECT 'c0000000-0000-0000-0000-000000000031', 'd0000000-0000-0000-0000-000000000001'
WHERE EXISTS (SELECT 1 FROM products WHERE id = 'c0000000-0000-0000-0000-000000000031')
  AND NOT EXISTS (SELECT 1 FROM product_tags WHERE product_id = 'c0000000-0000-0000-0000-000000000031' AND tag_id = 'd0000000-0000-0000-0000-000000000001');
