import ollama from 'ollama';
import openai from 'openai';
import { v4 as uuidv4 } from 'uuid';
import { getDescriptionsForSuspect } from './main';
import type { Answer, Question, Service, ServiceStatus } from './main';
import { get } from 'svelte/store';

import { activeService, services } from './stores';

const answerReflection = `ROLE: You are a player of Unusual Suspects board game - text based version. You are a witness.
TASK: Read the description of the perpetrator and the question the police officer asked you about perpetrator.
Write a short reflection on the perpetrator in relation to the question.
Try to think both ways, both about the positive answer and the negative one, which one you lean more towards. Cca 100 words.
QUESTION: %s
DESCRIPTION OF PERPETRATOR: %s`

const answerBoolean = `ROLE: You are a senior decision maker.
TASK: Answer the question YES or NO. Do not write anything else. Do not write anything else. Just write YES, or NO based on the previous information.`


export async function generateAnswer(roundUUID: string, question: Question, criminalUUID: string): Promise<Answer> {
    console.log(">>> getAnswerFromAI called!");    
    const service = get(services)[get(activeService)];

    try {    
        const descriptions = await getDescriptionsForSuspect(criminalUUID, service.Name, service.VisualModel);
        const description = descriptions.join("\n");
        
        let answer: string;
        const service_name = service.Name.toLowerCase();
        switch(service_name) {
            case "openai": 
                answer = await getAnswerFromOpenAI(question.English, description, service);
                break;
            case "ollama":
                answer = await getAnswerFromOllama(question.English, description, service);
                break;
            default:
                console.error(`Unsupported service '${service.Name}'`);
                await saveAnswer("failed: unsupported service", roundUUID);
                return;
        }

        const answer_object: Answer = {
            uuid: uuidv4(),
            answer: answer
        }
        console.log("Answer is:", answer_object);        
        return answer_object;
    } catch (error) {
        console.error(`getAnswerFromAI error for round ${roundUUID}:`, error);
    }
}


async function getAnswerFromOpenAI(question: string, description: string, service: Service): Promise<string> {
    return ""
}

async function getAnswerFromOllama(question: string, description: string, service: Service): Promise<string> {
    // First get reflection
    const reflectionPrompt = answerReflection.replace("%s", question).replace("%s", description);
    const reflectionResponse = await ollama.chat({
        model: service.VisualModel,
        messages: [{ role: "user", content: reflectionPrompt }]
    });
    const reflection = reflectionResponse.message.content;
    console.log("🤖 GOT THE REFLECTION FROM OLLAMA:", reflection);

    // Then get yes/no decision
    const decisionResponse = await ollama.chat({
        model: service.VisualModel,
        messages: [
            { role: "user", content: reflectionPrompt },
            { role: "assistant", content: reflection },
            { role: "user", content: answerBoolean }
        ]
    });
    const decision = decisionResponse.message.content;
    console.log("🤖 GOT THE FINAL ANSWER FROM OLLAMA:", decision);

    return decision;
}

async function saveAnswer(answer: string, roundUUID: string): Promise<void> {
    console.log("Saving answer:", answer);
}


// MARK: OLLAMA
export async function checkServiceStatusOllama(service: Service): Promise<ServiceStatus> {
    const status: ServiceStatus = {
        service: service,
        ready: false,
        message: "Ollama not ready"
    };

    try {
        const response = await ollama.list();
        if (!response) {
            status.message = "Ollama response is nil";
            return status;
        }

        for (const model of response.models) {
            if (model.name === service.VisualModel) {
                status.ready = true;
                break;
            }
            const [name] = model.name.split(":");
            if (name === service.VisualModel) {
                status.ready = true;
                break;
            }
        }

        status.message = "Ollama is running";
        return status;

    } catch (error) {
        status.message = `Ollama error: ${error}`;
        return status;
    }
}
