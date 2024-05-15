package osfile

import (
	"bufio"
	"encoding/json"
	"os"
)

type Event struct {
	UUID        string `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

type Producer struct {
	file   *os.File
	writer *bufio.Writer
}

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &Producer{file, bufio.NewWriter(file)}, nil
}

func (p *Producer) WriteEvent(event *Event) error {
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}
	// записываем событие в буфер
	if _, err := p.writer.Write(data); err != nil {
		return err
	}
	//добавляем перенос строки
	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}
	//записываем буфер в файл
	return p.writer.Flush()

}

func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &Consumer{file, bufio.NewScanner(file)}, nil
}

func (c *Consumer) ReadEvent() (*Event, error) {
	//одиночное сканирование до следующей строки
	if !c.scanner.Scan() {
		return nil, c.scanner.Err()
	}
	//читаем данные из сканнера
	data := c.scanner.Bytes()

	//преобразуем данные из JSON в структуру
	event := Event{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}
