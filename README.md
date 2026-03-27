## 🏛 Architecture Overview

This project demonstrates a modern FinTech workflow for credit card issuance, balancing rigid business rules (BPMN/DMN) with flexible AI decision-making (RAG).

### System Components:
* **Orchestration Layer:** [Camunda 8 SaaS](https://camunda.com/8/) (Zeebe Engine) managing the end-to-end process lifecycle.
* **Worker Layer:** High-performance microservices written in **Go (Golang)**, connecting via gRPC.
* **AI/RAG Engine:** * **LLM:** Google Vertex AI (Gemini 1.5 Flash) for automated risk assessment.
    * **Vector DB:** [Qdrant](https://qdrant.tech/) for storing and retrieving historical credit cases (Context Injection).
* **Cloud Infrastructure (GCP):**
    * **Compute:** Google Compute Engine (e2-micro) / Cloud Run.
    * **Database:** Google Firestore for application data persistence.
    * **IaC:** Terraform for reproducible infrastructure.

### The Pipeline Flow:
1. **Application Received:** Process starts via API or Camunda Form.
2. **Data Enrichment:** Go Worker fetches additional client data.
3. **Automated Scoring (DMN):** Fast-track rejection for non-eligible candidates.
4. **AI Risk Assessment (RAG):**
    * Worker generates embeddings for the current application.
    * Performs similarity search in **Qdrant** for past fraud/default patterns.
    * **Gemini** analyzes the combined data and provides a final recommendation.
5. **Decision Gateway:** Automated Approval or Rejection based on AI & Business logic.
