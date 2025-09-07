TARGET ?= studify
MODULE_NAME := $(shell grep '^module ' go.mod | cut -d' ' -f2)

BIN_DIR := ./.bin

SHELL := /bin/bash
GREEN  = \033[1;32m
YELLOW = \033[1;33m
PURPLE = \033[1;35m
RED    = \033[1;31m
RESET  = \033[0m

TOOLS = mvdan.cc/gofumpt@latest \
        github.com/daixiang0/gci@latest \
        github.com/segmentio/golines@latest \
        github.com/golangci/golangci-lint/cmd/golangci-lint@latest

COVERAGE_FILE ?= coverage.out

.PHONY: help
help: ## Показать список доступных для использования джоб
	@echo -e "$(PURPLE)Доступные джобы:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "$(YELLOW)%-20s$(RESET) %s\n", $$1, $$2}'

.PHONY: intstall-tools
intstall-tools: ## Установка тулзов для форматирования и линта кода (gofumpt, golines, gci, golangci-lint)
	@echo -e "$(GREEN)[INFO][DEPS-TOOLS][STARTED]$(RESET) Установка тулзов для форматирования и линта кода"
	@for tool in $(TOOLS); do \
		echo -e "$(PURPLE)  - Устанавливается $$tool$(RESET)"; \
		go install $$tool || { \
			echo -e "$(RED)[ERROR][DEPS-TOOLS][FAIL]$(RESET) Ошибка при установке $$tool"; exit 1; }; \
	done
	@echo -e "$(GREEN)[INFO][DEPS-TOOLS][SUCCESS]$(RESET) Все тулзы установлены"

.PHONY: build
build: ## Сборка бинарника приложения (через go build)
	@echo -e "$(YELLOW)[INFO][BUILD][STARTED]$(RESET) Выполняется go build для таргета $(TARGET)"
	@go build -o $(BIN_DIR)/$(TARGET) ./cmd/$(TARGET) && \
		echo -e "$(GREEN)[INFO][BUILD][SUCCESS]$(RESET) Сборка программы для таргета $(TARGET) завершена" || \
		echo -e "$(RED)[ERROR][BUILD][FAIL]$(RESET) Ошибка сборки программы для таргета $(TARGET)"
	@chmod +x $(BIN_DIR)/$(TARGET)

.PHONY: run
run: build ## Запуск программы (go build + go run)
	@echo -e "$(YELLOW)[INFO][RUN][STARTED]$(RESET) Запуск программы $(TARGET)"
	@$(BIN_DIR)/$(TARGET) && \
		echo -e "$(GREEN)[INFO][RUN][SUCCESS]$(RESET) Программа завершила работу" || \
		echo -e "$(RED)[ERROR][RUN][FAIL]$(RESET) Ошибка выполнения"

.PHONY: clean
clean: ## Очистка бинарников, кеша и остального мусора
	@echo -e "$(YELLOW)[INFO][CLEAN][STARTED]$(RESET) Очистка бинарников, кеша и мусора из $(BIN_DIR)"
	@(rm -rf $(BIN_DIR)/* && go clean) && \
		echo -e "$(GREEN)[INFO][CLEAN][SUCCESS]$(RESET) Очистка завершена" || \
		echo -e "$(RED)[ERROR][CLEAN][FAIL]$(RESET) Ошибка при очистке"

.PHONY: deps
deps: ## Установка зависимостей (go mod tidy + go mod download)
	@echo -e "$(YELLOW)[INFO][DEPS][STARTED]$(RESET) Установка зависимостей проекта"
	@(go mod tidy && go mod download) && \
		echo -e "$(GREEN)[INFO][DEPS][SUCCESS]$(RESET) Зависимости успешно установлены" || \
		echo -e "$(RED)[ERROR][DEPS][FAIL]$(RESET) Ошибка при установке зависимостей"

.PHONY: fmt-gofumpt
fmt-gofumpt: ## Форматирование кода через gofumpt
	@echo -e "$(YELLOW)[INFO][FMT-GOFUMPT][STARTED]$(RESET) Форматирование кода через gofumpt"
	@gofumpt -w . && \
		echo -e "$(GREEN)[INFO][FMT-GOFUMPT][SUCCESS]$(RESET) gofumpt завершил работу" || \
		echo -e "$(RED)[ERROR][FMT-GOFUMPT][FAIL]$(RESET) Ошибка при работе gofumpt"

.PHONY: fmt-golines
fmt-golines: ## Форматирование кода через golines
	@echo -e "$(YELLOW)[INFO][FMT-GOLINES][STARTED]$(RESET) Форматирование кода через golines"
	@golines -w . --max-len=100 && \
		echo -e "$(GREEN)[INFO][FMT-GOLINES][SUCCESS]$(RESET) golines завершил работу" || \
		echo -e "$(RED)[ERROR][FMT-GOLINES][FAIL]$(RESET) Ошибка при работе golines"

.PHONY: fmt-gci
fmt-gci: ## Форматирование кода через gci
	@echo -e "$(YELLOW)[INFO][FMT-GCI][STARTED]$(RESET) Форматирование кода через gci"
	@gci write --skip-generated -s standard -s default -s "prefix($(MODULE_NAME))" . && \
		echo -e "$(GREEN)[INFO][FMT-GCI][SUCCESS]$(RESET) gci завершил работу" || \
		echo -e "$(RED)[ERROR][FMT-GCI][FAIL]$(RESET) Ошибка при работе gci"

.PHONY: format
format: fmt-gofumpt fmt-golines fmt-gci ## Форматирование кода (gofumpt + golines + gci)
	@echo -e "$(YELLOW)[INFO][FORMAT][STARTED]$(RESET) Запуск форматирования кода"
	@$(MAKE) fmt-gofumpt && $(MAKE) fmt-golines && $(MAKE) fmt-gci && \
		echo -e "$(GREEN)[INFO][FORMAT][SUCCESS]$(RESET) Весь код успешно отформатирован" || \
		echo -e "$(RED)[ERROR][FORMAT][FAIL]$(RESET) Ошибка во время форматирования"

.PHONY: lint
lint: ## Линтинг кода через golangci-lint
	@echo -e "$(YELLOW)[INFO][LINT][STARTED]$(RESET) Запуск golangci-lint"
	@golangci-lint run ./... && \
		echo -e "$(GREEN)[INFO][LINT][SUCCESS]$(RESET) Линт прошел успешно" || \
		echo -e "$(RED)[ERROR][LINT][FAIL]$(RESET) Ошибки после линта"

.PHONY: test
test: ## Запуск тестов
	@echo -e "$(YELLOW)[INFO][TEST][STARTED]$(RESET) Запуск тестов"
	@go test -v ./... && \
		echo -e "$(GREEN)[INFO][TEST][SUCCESS]$(RESET) Все тесты прошли успешно" || \
		echo -e "$(RED)[ERROR][TEST][FAIL]$(RESET) Тесты завершились с ошибками"

.PHONY: coverage
coverage: ## Запуск тестов с проверкой покрытия
	@echo -e "$(YELLOW)[INFO][COVERAGE][STARTED]$(RESET) Запуск тестов с проверкой покрытия"
	@go test -coverpkg='$(MODULE_NAME)/...' --race -count=1 -coverprofile='$(COVERAGE_FILE)' ./... && \
		go tool cover -func='$(COVERAGE_FILE)' | grep ^total | tr -s '\t' && \
		echo -e "$(GREEN)[INFO][COVERAGE][SUCCESS]$(RESET) Тесты успешно прошли" || \
		echo -e "$(RED)[ERROR][COVERAGE][FAIL]$(RESET) Ошибка при работе тестов с проверкой покрытия"

.PHONY: bench
bench: ## Запуск бенчмарк тестов
	@echo -e "$(YELLOW)[INFO][BENCH][STARTED]$(RESET) Запуск бенчмарк тестов"
	@go test -bench=. ./... && \
		echo -e "$(GREEN)[INFO][BENCH][SUCCESS]$(RESET) Бенчмарки завершены" || \
		echo -e "$(RED)[ERROR][BENCH][FAIL]$(RESET) Ошибка при выполнении бенчмарков"

