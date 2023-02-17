/*
copyright 2020 the Escli authors

licensed under the apache license, version 2.0 (the "license");
you may not use this file except in compliance with the license.
you may obtain a copy of the license at

    http://www.apache.org/licenses/license-2.0

unless required by applicable law or agreed to in writing, software
distributed under the license is distributed on an "as is" basis,
without warranties or conditions of any kind, either express or implied.
see the license for the specific language governing permissions and
limitations under the license.
*/

package constants

import "os"

const (
	DefaultRegion  = "ap-northeast-2"
	EmptyString    = ""
	GlacierType    = "GLACIER"
	DefaultProduct = "elasticsearch"
	OpenSearch     = "opensearch"
	DefaultSortKey = "index"
)

const (
	GetIndexSetting         = 1
	GetIndexSettingWithName = 2
	PutIndexSetting         = 3

	GetClusterSetting = 0
	PutClusterSetting = 3

	HardLimitMaxConcurrentJob = 100
	DefaultMaxConcurrentJob   = 50
)

var (
	ConfigDirectoryPath     = HomeDir() + "/.escli"
	BaseFilePath            = ConfigDirectoryPath + "/config.yaml"
	ValidRestoreTier        = []string{"Standard", "Bulk", "Expedited"}
	SupportedRepositoryType = []string{"s3"}
)

// Get Home Directory
func HomeDir() string {
	if h := os.Getenv("HOME"); h != EmptyString {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
