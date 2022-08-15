/*
 * Copyright 2022 Sue B.V.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package redis

/* TODO redis mock doesn't support v9
var testContainerImage = model.ContainerImage{
	Hash:     "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
	HashType: "sha256",
	Status:   model.BeingFuzzed,
}

func createMockStorage() (client containerImageRedis, mock redismock.ClientMock) {
	db, mock := redismock.NewClientMock()
	repo := containerImageRedis{
		client: db,
	}

	return repo, mock
}

func TestCreateContainerImage(t *testing.T) {
	repo, mock := createMockStorage()

	expectedVal := strconv.Itoa(int(testContainerImage.Status))
	expectedKey := fmt.Sprintf("%s:%s", testContainerImage.HashType, testContainerImage.Hash)
	expectedExp := time.Duration(0)
	mock.ExpectSet(expectedKey, expectedVal, expectedExp).SetVal(expectedVal)
	err := repo.Create(context.TODO(), testContainerImage)
	assert.NoError(t, err)
}

func TestUpdateContainerImage(t *testing.T) {
	repo, mock := createMockStorage()

	expectedVal := strconv.Itoa(int(testContainerImage.Status))
	expectedKey := fmt.Sprintf("%s:%s", testContainerImage.HashType, testContainerImage.Hash)
	expectedExp := time.Duration(0)
	mock.ExpectSet(expectedKey, expectedVal, expectedExp).SetVal(expectedVal)
	err := repo.Update(context.TODO(), testContainerImage)
	assert.NoError(t, err)
}

func TestFindByHash(t *testing.T) {
	repo, mock := createMockStorage()

	strHash, strStatus := testContainerImage.String()

	// mock store returns an error if you don't add the ExpectSet line :(
	mock.ExpectSet(strHash, strStatus, time.Duration(0)).SetVal(strStatus)
	createErr := repo.client.Set(context.TODO(), strHash, strStatus, time.Duration(0)).Err()
	if !assert.NoError(t, createErr) {
		return
	}

	mock.ExpectGet(strHash).SetVal(strStatus)
	foundImage, didFind, findErr := repo.GetByKey(context.TODO(), strHash)

	assert.NoError(t, findErr)
	assert.True(t, didFind)
	assert.Equal(t, testContainerImage.Hash, foundImage.Hash)
	assert.Equal(t, testContainerImage.HashType, foundImage.HashType)
	assert.Equal(t, testContainerImage.Status, foundImage.Status)
}
*/

/* func TestCheckHealth(t *testing.T) {
	repo, mock := createMockStorage()

	mock.ExpectPing().SetVal("")
	pingResult := repo.CheckHealth(context.TODO())
	assert.True(t, pingResult.IsHealthy)
	assert.Equal(t, pingResult.Info[health.StatusKey], health.HealthyStatus, "status should be healthy and set to '"+health.HealthyStatus+"'")
}

func TestCheckHealthError(t *testing.T) {
	repo, mock := createMockStorage()

	errMsg := "some ping error"
	mock.ExpectPing().SetErr(fmt.Errorf(errMsg))
	pingResult := repo.CheckHealth(context.TODO())
	assert.False(t, pingResult.IsHealthy)
	assert.Equal(t, pingResult.Info[health.StatusKey], health.UnHealthyStatus, "status should be unhealthy and set to '"+health.UnHealthyStatus+"'")
	assert.Equal(t, pingResult.Info["reason"], errMsg, "reason should be set to error message '"+errMsg+"'")
} */
