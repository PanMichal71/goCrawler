package main

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockIStorage struct {
	mock.Mock
}

func (m *MockIStorage) Write(bytes []byte) error {
	args := m.Called(bytes)
	return args.Error(0)
}

func (m *MockIStorage) Open(filename string) error {
	args := m.Called(filename)
	return args.Error(0)
}

func (m *MockIStorage) Close() {
	m.Called()
}

type MockIDatabase struct {
	mock.Mock
}

func (m *MockIDatabase) Store(key string, value []byte) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockIDatabase) Read(key string) ([]byte, error) {
	args := m.Called(key)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockIDatabase) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockIDatabase) Exists(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m *MockIDatabase) ListKeys() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockIDatabase) Count() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

var defaultHtmlContent = "<html><body><a href=\"https://www.google.com\">Google</a></body></html>"
var defaultHtmlContentMd5Hash = "d6165a2f6a47eba8aa611ca6891203a9"

var changedHtmlContent = "<html><body><a href=\"https://www.google.com\">Yahoo</a></body></html>"
var changedHtmlContentMd5Hash = "2f2180839c2f324971d4f0f98fbf46de"

func fileNameMatchesPattern(version int) interface{} {
	return func(fileName string) bool {
		regexPattern := `^.+\/v` + fmt.Sprint(version) + `.html$`

		matched, err := regexp.MatchString(regexPattern, fileName)
		if err != nil {
			return false
		}

		return matched
	}
}

func Test_ShouldStoreUrlVersionsInSeparateDirectories(t *testing.T) {
	storageMock := new(MockIStorage)
	databaseMock := new(MockIDatabase)

	sut := NewDifferenceTracker(databaseMock, storageMock)

	storageMock.On("Open", mock.MatchedBy(fileNameMatchesPattern(1))).Return(nil).Once()
	storageMock.On("Write", []byte(defaultHtmlContent)).Return(nil).Once()
	storageMock.On("Close").Return().Once()

	databaseMock.On("Exists", "https://www.google.com").Return(false, nil).Once()
	jsonWithSingleVersion := []byte(`[{"Hash":"d6165a2f6a47eba8aa611ca6891203a9","FilePath":"google.com/v1.html","Version":1}]`)

	databaseMock.On("Store", "https://www.google.com", jsonWithSingleVersion).Return(nil).Once()

	sut.HandleContent("https://www.google.com", defaultHtmlContent)

	storageMock.On("Open", mock.MatchedBy(fileNameMatchesPattern(2))).Return(nil)
	storageMock.On("Write", []byte(changedHtmlContent)).Return(nil)
	storageMock.On("Close").Return().Once()

	databaseMock.On("Exists", "https://www.google.com").Return(true, nil)
	databaseMock.On("Read", "https://www.google.com").Return(jsonWithSingleVersion, nil)

	jsonWithTwoVersions := []byte(`[{"Hash":"d6165a2f6a47eba8aa611ca6891203a9","FilePath":"google.com/v1.html","Version":1},{"Hash":"2f2180839c2f324971d4f0f98fbf46de","FilePath":"google.com/v2.html","Version":2}]`)
	databaseMock.On("Store", "https://www.google.com", jsonWithTwoVersions).Return(nil)

	sut.HandleContent("https://www.google.com", changedHtmlContent)

	storageMock.AssertNumberOfCalls(t, "Close", 2)
	storageMock.AssertNumberOfCalls(t, "Open", 2)
	storageMock.AssertNumberOfCalls(t, "Write", 2)
	databaseMock.AssertNumberOfCalls(t, "Exists", 2)
	databaseMock.AssertNumberOfCalls(t, "Read", 1)
	databaseMock.AssertNumberOfCalls(t, "Store", 2)
}

func Test_ShouldNotStoreSameVersionTwice(t *testing.T) {
	storageMock := new(MockIStorage)
	databaseMock := new(MockIDatabase)

	sut := NewDifferenceTracker(databaseMock, storageMock)

	storageMock.On("Open", mock.MatchedBy(fileNameMatchesPattern(1))).Return(nil).Once()
	storageMock.On("Write", []byte(defaultHtmlContent)).Return(nil).Once()
	storageMock.On("Close").Return().Once()

	databaseMock.On("Exists", "https://www.google.com").Return(false, nil).Once()
	jsonWithSingleVersion := []byte(`[{"Hash":"d6165a2f6a47eba8aa611ca6891203a9","FilePath":"google.com/v1.html","Version":1}]`)

	databaseMock.On("Store", "https://www.google.com", jsonWithSingleVersion).Return(nil).Once()

	sut.HandleContent("https://www.google.com", defaultHtmlContent)

	storageMock.On("Open", mock.MatchedBy(fileNameMatchesPattern(2))).Return(nil)
	storageMock.On("Write", []byte(defaultHtmlContent)).Return(nil)
	storageMock.On("Close").Return().Once()

	databaseMock.On("Exists", "https://www.google.com").Return(true, nil)
	databaseMock.On("Read", "https://www.google.com").Return(jsonWithSingleVersion, nil)

	sut.HandleContent("https://www.google.com", defaultHtmlContent)

	storageMock.AssertNumberOfCalls(t, "Close", 1)
	storageMock.AssertNumberOfCalls(t, "Open", 1)
	storageMock.AssertNumberOfCalls(t, "Write", 1)
	databaseMock.AssertNumberOfCalls(t, "Exists", 2)
	databaseMock.AssertNumberOfCalls(t, "Read", 1)
	databaseMock.AssertNumberOfCalls(t, "Store", 1)
}

func Test_ShouldPanicIfErrorOccursWhenOpeningStorage(t *testing.T) {
	storageMock := new(MockIStorage)
	databaseMock := new(MockIDatabase)

	sut := NewDifferenceTracker(databaseMock, storageMock)

	storageMock.On("Open", mock.MatchedBy(fileNameMatchesPattern(1))).Return(fmt.Errorf("error")).Once()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	sut.HandleContent("https://www.google.com", defaultHtmlContent)
}

func Test_ShouldStoreSeparateDomains(t *testing.T) {
	storageMock := new(MockIStorage)
	databaseMock := new(MockIDatabase)

	sut := NewDifferenceTracker(databaseMock, storageMock)

	storageMock.On("Open", mock.MatchedBy(fileNameMatchesPattern(1))).Return(nil).Once()
	storageMock.On("Write", []byte(defaultHtmlContent)).Return(nil).Once()
	storageMock.On("Close").Return().Once()

	databaseMock.On("Exists", "https://www.google.com").Return(false, nil).Once()
	jsonWithSingleVersion := []byte(`[{"Hash":"d6165a2f6a47eba8aa611ca6891203a9","FilePath":"google.com/v1.html","Version":1}]`)

	databaseMock.On("Store", "https://www.google.com", jsonWithSingleVersion).Return(nil).Once()

	sut.HandleContent("https://www.google.com", defaultHtmlContent)

	storageMock.On("Open", mock.MatchedBy(fileNameMatchesPattern(1))).Return(nil)
	storageMock.On("Write", []byte(changedHtmlContent)).Return(nil)
	storageMock.On("Close").Return().Once()

	databaseMock.On("Exists", "https://www.google2.com").Return(false, nil)

	jsonWithTwoVersions := []byte(`[{"Hash":"2f2180839c2f324971d4f0f98fbf46de","FilePath":"google2.com/v1.html","Version":1}]`)
	databaseMock.On("Store", "https://www.google2.com", jsonWithTwoVersions).Return(nil)

	sut.HandleContent("https://www.google2.com", changedHtmlContent)

	storageMock.AssertNumberOfCalls(t, "Close", 2)
	storageMock.AssertNumberOfCalls(t, "Open", 2)
	storageMock.AssertNumberOfCalls(t, "Write", 2)
	databaseMock.AssertNumberOfCalls(t, "Exists", 2)
	databaseMock.AssertNumberOfCalls(t, "Read", 0)
	databaseMock.AssertNumberOfCalls(t, "Store", 2)
}
