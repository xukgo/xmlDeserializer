/*
@Time : 2019/11/11 16:16
@Author : Hermes
@File : checkSliceContains_test
@Description:
*/
package xmlDeserializer

import "testing"

func TestCeckSliceContains(t *testing.T) {
	var arr = []string{"11", "22"}
	dest := "11"

	actual := checkStringSliceContains(arr, dest)
	if !actual {
		t.Fail()
	}
}

func BenchmarkSliceContains(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var arr = []string{"11", "22"}
		dest := "11"

		actual := checkStringSliceContains(arr, dest)
		if !actual {
			b.Fail()
		}
	}
}
