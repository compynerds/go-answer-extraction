package main

import (
	"encoding/json"
	"fmt"
	"github.com/nuvi/go-healthz"
	"github.com/nuvi/go-models"
	"github.com/nuvi/go-rabbitmq"
	"github.com/nuvi/unicycle"
	"log"
)

const (
	fromQueuePrefix    = "answer-extraction-q"
	fromRoutingKey     = "tagger.tagged_compact_engage_activity.created"
	toSurveyRoutingKey = "answer_extraction.survey_response_activity.created"
	toAnswerRoutingKey = "answer_extraction.survey_answer_activity.created"
	concurrency        = 10
)

func main() {
	hapi := healthz.StartDefaultHealthChecks()

	rabbitURL := unicycle.GetenvOrFatal("RABBITMQ_URL")

	consumeConn, err := rabbitmq.NewConsumer(
		rabbitURL,
		rabbitmq.Config{},
		rabbitmq.WithConsumerOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer consumeConn.Close()

	publisher, err := rabbitmq.NewPublisher(rabbitURL, rabbitmq.Config{}, rabbitmq.WithPublisherOptionsLogging)
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()
	returns := publisher.NotifyReturn()
	go func() {
		for r := range returns {
			log.Printf("message returned from rabbit server: %s", string(r.Body))
		}
	}()

	hapi.SetReadiness(true)

	handler := getSubscriberHandler(publisher)

	err = consumeConn.StartConsuming(
		handler,
		fmt.Sprintf("%s.%s", fromQueuePrefix, fromRoutingKey),
		[]string{fromRoutingKey},
		rabbitmq.WithConsumeOptionsConcurrency(concurrency),
		rabbitmq.WithConsumeOptionsQOSPrefetch(concurrency*2),
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsQuorum,
		rabbitmq.WithConsumeOptionsBindingExchangeName("events"),
		rabbitmq.WithConsumeOptionsBindingExchangeKind("topic"),
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)
	<-forever
}

func getSubscriberHandler(publisher *rabbitmq.Publisher) func(d rabbitmq.Delivery) rabbitmq.Action {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		compactSurvey := models.CompactSurveyActivity{}
		err := json.Unmarshal(d.Body, &compactSurvey)
		if err != nil {
			log.Printf("error could not unmarshal compact survey with: %v", err)
			return rabbitmq.NackRequeue
		}
		surveyResponse := models.CompactSurveyResponseActivity{
			CompactActivity:         compactSurvey.CompactActivity,
			SurveyTemplateID:        compactSurvey.SurveyTemplateID,
			SurveyResultID:          compactSurvey.SurveyResultID,
			SurveySentAt:            compactSurvey.SurveySentAt,
			SurveyRespondedAt:       compactSurvey.SurveyRespondedAt,
			SurveyResponseUpdatedAt: compactSurvey.SurveyResponseUpdatedAt,
			SurveyResponseOrigin:    compactSurvey.SurveyResponseOrigin,
			SurveyFinalTakenAt:      compactSurvey.SurveyFinalTakenAt,
			CustomerID:              compactSurvey.CustomerID,
			RepSentimentCategory:    compactSurvey.RepSentimentCategory,
			MetaData:                compactSurvey.MetaData,
			SurveyCategories:        compactSurvey.SurveyCategories,
			KeywordIDs:              compactSurvey.KeywordIDs,
			CustomScores:            compactSurvey.CustomScores,
		}
		answers := []models.CompactSurveyAnswerActivity{}

		for _, answer := range compactSurvey.SurveyAnswers {
			answers = append(answers, models.CompactSurveyAnswerActivity{
				CompactActivity:         compactSurvey.CompactActivity,
				SurveyTemplateId:        compactSurvey.SurveyTemplateID,
				SurveyResultId:          compactSurvey.SurveyResultID,
				SurveyQuestionId:        answer.SurveyQuestionId,
				SurveyQuestionTextEn:    answer.SurveyQuestionTextEn,
				ParentQuestionTextEn:    answer.ParentQuestionTextEn,
				ParentQuestionId:        answer.ParentQuestionId,
				SurveySentAt:            answer.SurveySentAt,
				SurveyRespondedAt:       compactSurvey.SurveyRespondedAt,
				SurveyResponseUpdatedAt: compactSurvey.SurveyResponseUpdatedAt,
				//SurveyResponseOrigin: compactSurvey.SurveyResponseOrigin, // TODO: what type is this and is this even the correct field name?
				CustomerId:                  compactSurvey.CustomerID,
				SurveyTakerName:             answer.SurveyTakerName,
				SurveyTakerPhone:            answer.SurveyTakerPhone,
				SurveyTakerEmail:            answer.SurveyTakerEmail,
				SurveyQuestionType:          answer.SurveyQuestionType,
				SurveyAnswerType:            answer.SurveyAnswerType,
				SurveyAnswerText:            answer.SurveyAnswerText,
				SurveyAnswerKeyword:         answer.SurveyAnswerKeyword,
				SurveyAnswerBool:            answer.SurveyAnswerBool,
				SurveyAnswerNumber:          answer.SurveyAnswerNumber,
				SurveyAnswerDate:            answer.SurveyAnswerDate,
				SurveyAnswerRating:          answer.SurveyAnswerRating,
				SurveyAnswerOptionId:        answer.SurveyAnswerOptionId,
				SurveyAnswerOptionEn:        answer.SurveyAnswerOptionEn,
				SurveyAnswerMinRating:       answer.SurveyAnswerMinRating,
				SurveyAnswerMaxRating:       answer.SurveyAnswerMaxRating,
				SurveyTopBoxRank:            answer.SurveyTopBoxRank,
				SurveyBottomBoxRank:         answer.SurveyBottomBoxRank,
				SurveyQuestionSkipped:       answer.SurveyQuestionSkipped,
				RawBodyTextEncryptedSfw:     compactSurvey.RawBodyText, // TODO: change this to the encrypted version
				RawBodyTextEncrypted:        compactSurvey.RawBodyText, // TODO: change this to the encrypted version
				ReputationSentimentScore:    answer.ReputationSentimentScore,
				ReputationSentimentCategory: answer.ReputationSentimentCategory,
				Metadata:                    answer.Metadata,
				IsTopBox:                    answer.IsTopBox,
				TimeToRespond:               answer.TimeToRespond,
				SurveyQuestionSeenAt:        answer.SurveyQuestionSeenAt,
				RespondedAt:                 answer.RespondedAt,
				IsPii:                       answer.IsPii,
				VisibleToUserRoles:          answer.VisibleToUserRoles,
			})
		}
		respData, respMarshErr := json.Marshal(surveyResponse)
		if respMarshErr != nil {
			log.Printf("error marshalling response with: %v nacking", respMarshErr)
			return rabbitmq.NackRequeue
		}

		publisher.Publish(
			respData, // if they send us bogus json this could break but since this is a temp measure we'll try this
			[]string{toSurveyRoutingKey},
			rabbitmq.WithPublishOptionsContentType("application/json"),
			rabbitmq.WithPublishOptionsPersistentDelivery,
			rabbitmq.WithPublishOptionsExchange("events"),
		)

		for _, answer := range answers {
			data, marshErr := json.Marshal(answer)
			if marshErr != nil {
				log.Printf("error marshalling answer with: %v", marshErr)
			}
			publisher.Publish(
				data, // if they send us bogus json this could break but since this is a temp measure we'll try this
				[]string{toAnswerRoutingKey},
				rabbitmq.WithPublishOptionsContentType("application/json"),
				rabbitmq.WithPublishOptionsPersistentDelivery,
				rabbitmq.WithPublishOptionsExchange("events"),
			)
		}
		return rabbitmq.Ack
	}
}
