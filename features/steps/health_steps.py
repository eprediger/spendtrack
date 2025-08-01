from behave import given, when, then, use_step_matcher

use_step_matcher("re")


@given("the app is running")
def step_impl(context):
    """
    :type context: behave.runner.Context
    """
    pass


@when("I request the health status")
def step_impl(context):
    """
    :type context: behave.runner.Context
    """
    context.response = context.client.get("/health")


@then("I should receive a healthy status")
def step_impl(context):
    """
    :type context: behave.runner.Context
    """
    assert context.response.status_code == 200
