resource "google_pubsub_topic" "plan_generated" {
  name = "plan.generated"
}

resource "google_pubsub_topic" "order_requested" {
  name = "order.requested"
}

resource "google_pubsub_topic" "feedback_logged" {
  name = "feedback.logged"
}
