echo ""
echo "🚀 Полное тестирование Notes API"
echo "==============================="


# Тест 1: Создание новой записи
echo ""
echo "🔍 Создание новой записи"
echo "Запрос: POST $BASE_URL/$SERVICE_NAME_NOTES/note"
echo "Ответ:"
CREATE_RESPONSE=$(curl -X "POST" "$BASE_URL/$SERVICE_NAME_NOTES/note" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"name":"Test Note","content":"Test Content"}' \
     -w "\n📊 HTTP Статус: %{http_code}\n")
     
# Извлекаем ID созданной записи из ответа
ID_NOTE=$(echo "$CREATE_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "ID созданной записи: $ID_NOTE"
echo "-------------------------------------------"

# Небольшая пауза между запросами
sleep 2
# Тест 2: Получение списка всех записей
echo ""
echo "🔍 Получение списка всех записей"
echo "Запрос: GET $BASE_URL/$SERVICE_NAME_NOTES/notes"
echo "Ответ:"
curl -X "GET" "$BASE_URL/$SERVICE_NAME_NOTES/notes" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -w "\n📊 HTTP Статус: %{http_code}\n"
echo "-------------------------------------------"

# Небольшая пауза между запросами
sleep 2
# Тест 3: Получение записи по ID
echo ""
echo "🔍 Получение записи по ID"
echo "Запрос: GET $BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE"
echo "Ответ:"
curl -X "GET" "$BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -w "\n📊 HTTP Статус: %{http_code}\n"
echo "-------------------------------------------"

# Небольшая пауза между запросами
sleep 2
# Тест 4: Редактирование записи по ID
echo ""
echo "🔍 Редактирование записи по ID"
echo "Запрос: PUT $BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE"
echo "Ответ:"
curl -X "PUT" "$BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"name":"Updated Note","content":"Updated Content"}' \
     -w "\n📊 HTTP Статус: %{http_code}\n"
echo "-------------------------------------------"

# Небольшая пауза между запросами
sleep 2
# Тест 5: Удаление записи по ID
echo ""
echo "🔍 Удаление записи по ID"
echo "Запрос: DELETE $BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE"
echo "Ответ:"
curl -X "DELETE" "$BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -w "\n📊 HTTP Статус: %{http_code}\n"
echo "-------------------------------------------"

echo "✅ Все тесты завершены!"
echo "==============================="