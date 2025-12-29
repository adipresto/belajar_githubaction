Feature: User manage their activities
  As User
  I can post new Activity
  Time pass by, somehow I record about 100 activities
  I can open my activities and filter by calories burned and activity type

  Scenario: I can post new Activity
    Given tokenFromRegister exists
    When I prepare a POST request to "/v1/activity"
    And I set a bearer authentication header with tokenFromRegister value
    And I use activity payload
    And send
    Then the response should be 201
    And the response body matches the activity schema
    And "caloriesBurned" value is equal to 150
    And save the "activityId" value as savedActivityId
    And doneAt must be the result of adding the durationInMinutes to createdAt

  Scenario: Overtime I post about 100 activities total
    Given tokenFromRegister exists
    When I prepare a POST request to "/v1/activity"
    And I set a bearer authentication header with tokenFromRegister value
    And I prepare about 100 activity to inject and send right away

  Scenario: I open my activity dashboard, fill with 5 datas
    Given tokenFromRegister exists
    When I prepare a GET request to "/v1/activity?limit=5&offset=0"
    And I set a bearer authentication header with tokenFromRegister value
    And send
    Then the response should be 200
    And the response body matches the activities schema
    And the fetched data total are 5 datas

  Scenario: I put filter about 300 to 500 calories burned
    Given tokenFromRegister exists
    When I prepare a GET request to "/v1/activity?limit=5&offset=0&caloriesBurnedMin=300&caloriesBurnedMax=500"
    And I set a bearer authentication header with tokenFromRegister value
    And send
    Then the response should be 200
    And the response body matches the activities schema
    And the value of caloriesBurned are between 300 and 500

  Scenario: Get five of activities type Yoga to PATCH, DELETE
    Given tokenFromRegister exists
    When I prepare a GET request to "/v1/activity?limit=5&offset=0&activityType=Yoga"
    And I set a bearer authentication header with tokenFromRegister value
    And send
    Then the response should be 200
    And the response body matches the activities schema
    And the value of activityType is Yoga
    And save the "0.activityId" value as targetActivityIdToUpdate
    And save the "1.activityId" value as targetActivityIdToDelete
    And save the "2.activityId" value as targetActivityIdToDeleteFromUnkownAccount
    And save the "3.activityId" value as targetActivityIdToDeleteFromOtherAccount

    Given targetActivityIdToUpdate exists
    When I prepare a PATCH request to "/v1/activity/" with saved value targetActivityIdToUpdate
    And I set a bearer authentication header with tokenFromRegister value
    And I use updateActivity payload
    And send
    Then the response should be 201
    And the response body matches the activity schema
    And "activityType" value is equal to "Running"
    And "caloriesBurned" value is equal to 450
    And doneAt must be the result of adding the durationInMinutes to createdAt

    Given targetActivityIdToDelete exists
    When I prepare a DELETE request to "/v1/activity/" with saved value targetActivityIdToDelete
    And I set a bearer authentication header with tokenFromRegister value
    And send
    Then the response should be 200

  # negative scenarios

  # POST
  Scenario: I didn't have account, I should not able to post activity
    When I prepare a POST request to "/v1/activity"
    And I set a bearer authentication header with false token
    And I use activity payload
    And send
    Then the response should be 401

  # DELETE
  Scenario: I didn't have account, I supposed to not able to delete
    When I prepare a DELETE request to "/v1/activity/" with saved value targetActivityIdToDeleteFromUnkownAccount
    And I set a bearer authentication header with false token
    And I use activity payload
    And send
    Then the response should be 401

  # PATCH
  Scenario: From other account, I supposed to not able to delete
    Given tokenNegativeScenario exists
    When I prepare a PATCH request to "/v1/activity/" with saved value targetActivityIdToDeleteFromOtherAccount
    And I set a bearer authentication header with tokenNegativeScenario value
    And I use activity payload
    And send
    Then the response should be 401

  Scenario: With other account, I must not able to update activity isn't mine
    Given tokenNegativeScenario exists
    And targetActivityIdToUpdate exists
    When I prepare a PATCH request to "/v1/activity/" with saved value targetActivityIdToUpdate
    And I set a bearer authentication header with tokenNegativeScenario value
    And I use updateActivity payload
    And send
    Then the response should be 404

  Scenario: I didn't have account, but I supposed to not able to delete
    Given tokenNegativeScenario exists
    When I prepare a PATCH request to "/v1/activity/" with saved value targetActivityIdToDeleteFromOtherAccount
    And I set a bearer authentication header with tokenNegativeScenario value
    And I use activity payload
    And send
    Then the response should be 401

  # GET
  Scenario: With other account, the activity should be counted one since I only add once to this account
    Given tokenNegativeScenario exists
    When I prepare a POST request to "/v1/activity"
    And I set a bearer authentication header with tokenNegativeScenario value
    And I use activity payload
    And send
    Then the response should be 201

    When I prepare a GET request to "/v1/activity?limit=5&offset=0"
    And I set a bearer authentication header with tokenNegativeScenario value
    And send
    Then the response body matches the activities schema
    And the fetched data total are 1 datas

