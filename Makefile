# Коммиты
.PHONY: commit
commit:
	@echo "Выберите тип коммита:"
	@echo "1. feat: новая функциональность"
	@echo "2. fix: исправление ошибки"
	@echo "3. docs: изменения в документации"
	@echo "4. style: форматирование кода"
	@echo "5. refactor: рефакторинг"
	@echo "6. test: добавление тестов"
	@echo "7. chore: обновление зависимостей, конфигурации, структуры проекта"
	@read -p "Выберите номер (1-7): " type_num; \
	case $$type_num in \
		1) type="feat";; \
		2) type="fix";; \
		3) type="docs";; \
		4) type="style";; \
		5) type="refactor";; \
		6) type="test";; \
		7) type="chore";; \
		*) echo "Неверный выбор"; exit 1;; \
	esac; \
	read -p "Введите область изменений (опционально, например: auth, user, api): " scope; \
	read -p "Введите описание коммита: " description; \
	if [ -n "$$scope" ]; then \
		git commit -m "$$type($$scope): $$description"; \
		echo "Коммит создан: $$type($$scope): $$description"; \
	else \
		git commit -m "$$type: $$description"; \
		echo "Коммит создан: $$type: $$description"; \
	fi