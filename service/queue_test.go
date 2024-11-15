package service

import (
	"errors"
	"testing"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/mocks"
	"github.com/stretchr/testify/assert"
)

func TestQueueService_Publish_Success(t *testing.T) {
	rabbitMQMock := new(mocks.RabbitMQClient)
	queueService := &queueService{
		rabbitMQClient: rabbitMQMock,
	}

	queueName := domain.QueueLikePost
	message := []byte("test message")

	rabbitMQMock.On("Publish", queueName, message).Return(nil)

	err := queueService.Publish(queueName, message)

	assert.NoError(t, err)
	rabbitMQMock.AssertExpectations(t)
}

func TestQueueService_Publish_Failure(t *testing.T) {
	rabbitMQMock := new(mocks.RabbitMQClient)
	queueService := &queueService{
		rabbitMQClient: rabbitMQMock,
	}

	queueName := domain.QueueLikePost
	message := []byte("test message")

	rabbitMQMock.On("Publish", queueName, message).Return(errors.New("publish error"))

	err := queueService.Publish(queueName, message)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error publishing message to queue")
	rabbitMQMock.AssertExpectations(t)
}

func TestQueueService_Consume_Success(t *testing.T) {
	rabbitMQMock := new(mocks.RabbitMQClient)
	queueService := &queueService{
		rabbitMQClient: rabbitMQMock,
	}

	queueName := "test_queue"
	messageChan := make(chan []byte, 1)
	messageChan <- []byte("test message")
	close(messageChan)

	var readOnlyMessageChan <-chan []byte = messageChan

	rabbitMQMock.On("Consume", queueName).Return(readOnlyMessageChan, nil)

	messages, err := queueService.Consume(queueName)

	assert.NoError(t, err)
	assert.Equal(t, readOnlyMessageChan, messages)
	rabbitMQMock.AssertExpectations(t)
}

func TestQueueService_Consume_Failure(t *testing.T) {
	rabbitMQMock := new(mocks.RabbitMQClient)
	queueService := &queueService{
		rabbitMQClient: rabbitMQMock,
	}

	queueName := "test_queue"

	rabbitMQMock.On("Consume", queueName).Return(nil, errors.New("consume error"))

	messages, err := queueService.Consume(queueName)

	assert.Error(t, err)
	assert.Nil(t, messages)
	assert.Contains(t, err.Error(), "error consuming message from queue")
	rabbitMQMock.AssertExpectations(t)
}
