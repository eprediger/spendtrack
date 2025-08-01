from behave import fixture, use_fixture
from fastapi.testclient import TestClient

from app.infrastructure.http.main import app

@fixture
def app_client(context, *args, **kwargs):
    context.client = TestClient(app)
    yield context.client


def before_all(context):
    pass


def before_scenario(context, scenario):
    if "todo" in scenario.effective_tags:
        scenario.skip("Scenario marked as @todo")


def before_feature(context, feature):
    if "todo" in feature.tags:
        feature.skip("Scenario marked as @todo")
        return

    use_fixture(app_client, context)
    context.vars = {}  # vars allow share variable between steps, reset them between features


def after_scenario(context, scenario):
    pass


def after_all(context):
    pass
