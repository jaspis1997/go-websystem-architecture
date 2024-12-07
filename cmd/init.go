package main

import (
	"playground/internal/entities"
	"playground/internal/external/log"
	"playground/internal/logics"
)

var (
	// ロギングライブラリの初期化 (optional)
	initLogger func() error
	// 環境変数のロード (optional)
	loadEnvironment func() error
	// バッチ処理の初期化 (optional)
	initBatch func() error

	// リポジトリの初期化 (required)
	initRepository func() (repo *entities.Repository, err error)
	// ロジックの初期化 (required)
	initLogics func(repo *entities.Repository) error
	// Webエンジンの初期化 (required)
	initWebEngine func() error
)

func checkRequiredFunctions() {
	if initRepository == nil {
		log.Fatal("required field is empty: initRepository")
	}
	if initLogics == nil {
		log.Fatal("required field is empty: initLogics")
	}
	if initWebEngine == nil {
		log.Fatal("required field is empty: initWebEngine")
	}
}

func initialize() (err error) {
	checkRequiredFunctions()
	if initLogger != nil {
		if err = initLogger(); err != nil {
			if err, ok := err.(*entities.StackTraceableError); ok {
				return err.Add(entities.ErrorInitializeFailed)
			}
			return entities.NewStackTrace().
				SetIdAutomatically().
				Add(err).
				Add(entities.ErrorInitializeFailed)
		}
	}
	if loadEnvironment != nil {
		if err = loadEnvironment(); err != nil {
			if err, ok := err.(*entities.StackTraceableError); ok {
				return err.Add(entities.ErrorInitializeFailed)
			}
			return entities.NewStackTrace().
				SetIdAutomatically().
				Add(err).
				Add(entities.ErrorInitializeFailed)
		}
	}
	repository, err = initRepository()
	if err != nil {
		if err, ok := err.(*entities.StackTraceableError); ok {
			return err.Add(entities.ErrorInitializeFailed)
		}
		return entities.NewStackTrace().
			SetIdAutomatically().
			Add(err).
			Add(entities.ErrorInitializeFailed)

	}
	if initBatch != nil {
		if err := initBatch(); err != nil {
			if err, ok := err.(*entities.StackTraceableError); ok {
				return err.Add(entities.ErrorInitializeFailed)
			}
			return entities.NewStackTrace().
				SetIdAutomatically().
				Add(err).
				Add(entities.ErrorInitializeFailed)
		}
	}
	if err := logics.Initialize(repository); err != nil {
		if err, ok := err.(*entities.StackTraceableError); ok {
			return err.Add(entities.ErrorInitializeFailed)
		}
		return entities.NewStackTrace().
			SetIdAutomatically().
			Add(err).
			Add(entities.ErrorInitializeFailed)

	}
	if err := initWebEngine(); err != nil {
		if err, ok := err.(*entities.StackTraceableError); ok {
			return err.Add(entities.ErrorInitializeFailed)
		}
		return entities.NewStackTrace().
			SetIdAutomatically().
			Add(err).
			Add(entities.ErrorInitializeFailed)
	}
	return nil
}
