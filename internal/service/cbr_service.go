package service

import (
    "bytes"
    "errors"
    "fmt"
    "io"
    "net/http"
    "time"
    
    "github.com/beevik/etree"
    "bank-service/internal/config"
)

type CBRService struct {
    config *config.Config
}

func NewCBRService(config *config.Config) *CBRService {
    return &CBRService{config: config}
}

func (s *CBRService) GetKeyRate() (float64, error) {
    soapRequest := s.buildSOAPRequest()
    rawBody, err := s.sendRequest(soapRequest)
    if err != nil {
        return 0, err
    }
    
    rate, err := s.parseXMLResponse(rawBody)
    if err != nil {
        return 0, err
    }
    
    // Add bank margin
    rate += 5.0
    return rate, nil
}

func (s *CBRService) buildSOAPRequest() string {
    fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
    toDate := time.Now().Format("2006-01-02")
    
    return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
        <soap12:Envelope xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
            <soap12:Body>
                <KeyRate xmlns="http://web.cbr.ru/">
                    <fromDate>%s</fromDate>
                    <ToDate>%s</ToDate>
                </KeyRate>
            </soap12:Body>
        </soap12:Envelope>`, fromDate, toDate)
}

func (s *CBRService) sendRequest(soapRequest string) ([]byte, error) {
    client := &http.Client{Timeout: 10 * time.Second}
    
    req, err := http.NewRequest(
        "POST",
        s.config.CentralBankURL,
        bytes.NewBuffer([]byte(soapRequest)),
    )
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
    req.Header.Set("SOAPAction", "http://web.cbr.ru/KeyRate")
    
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error sending request: %v", err)
    }
    defer resp.Body.Close()
    
    rawBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response: %v", err)
    }
    
    return rawBody, nil
}

func (s *CBRService) parseXMLResponse(rawBody []byte) (float64, error) {
    doc := etree.NewDocument()
    if err := doc.ReadFromBytes(rawBody); err != nil {
        return 0, fmt.Errorf("error parsing XML: %v", err)
    }
    
    krElements := doc.FindElements("//diffgram/KeyRate/KR")
    if len(krElements) == 0 {
        return 0, errors.New("key rate data not found")
    }
    
    latestKR := krElements[0]
    rateElement := latestKR.FindElement("./Rate")
    if rateElement == nil {
        return 0, errors.New("Rate tag not found")
    }
    
    var rate float64
    if _, err := fmt.Sscanf(rateElement.Text(), "%f", &rate); err != nil {
        return 0, fmt.Errorf("error converting rate: %v", err)
    }
    
    return rate, nil
}
