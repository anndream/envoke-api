package api

import (
	"testing"

	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/crypto/ed25519"
	ld "github.com/zbo14/envoke/linked_data"
	"github.com/zbo14/envoke/spec"
)

var CHALLENGE = "abc"

func GetPrivateKey(data Data) crypto.PrivateKey {
	privkey := new(ed25519.PrivateKey)
	privkey.FromString(data.GetStr("privateKey"))
	return privkey
}

func GetUserId(data Data) string {
	return data.GetStr("userId")
}

func TestApi(t *testing.T) {
	api := NewApi()
	output := MustOpenWriteFile("output.json")
	composer, err := api.Register(
		"itisasecret",
		spec.NewUser("composer@email.com", "", "", nil, "composer", "", "www.composer.com", "Person"),
	)
	if err != nil {
		t.Fatal(err)
	}
	composerId := GetUserId(composer)
	composerPrivkey := GetPrivateKey(composer)
	WriteJSON(output, composer)
	recordLabel, err := api.Register(
		"shhhh",
		spec.NewUser("record_label@email.com", "", "", nil, "record_label", "", "www.record_label.com", "Organization"),
	)
	if err != nil {
		t.Fatal(err)
	}
	recordLabelId := GetUserId(recordLabel)
	recordLabelPrivkey := GetPrivateKey(recordLabel)
	WriteJSON(output, recordLabel)
	performer, err := api.Register(
		"makeitup",
		spec.NewUser("performer@email.com", "123456789", "", nil, "performer", "ASCAP", "www.performer.com", "MusicGroup"),
	)
	if err != nil {
		t.Fatal(err)
	}
	performerId := GetUserId(performer)
	performerPrivkey := GetPrivateKey(performer)
	WriteJSON(output, performer)
	producer, err := api.Register(
		"1234",
		spec.NewUser("producer@email.com", "", "", nil, "producer", "", "www.soundcloud_page.com", "Person"),
	)
	if err != nil {
		t.Fatal(err)
	}
	producerId := GetUserId(producer)
	producerPrivkey := GetPrivateKey(producer)
	WriteJSON(output, producer)
	publisher, err := api.Register(
		"didyousaysomething?",
		spec.NewUser("publisher@email.com", "", "", nil, "publisher", "", "www.publisher.com", "Organization"),
	)
	if err != nil {
		t.Fatal(err)
	}
	publisherId := GetUserId(publisher)
	publisherPrivkey := GetPrivateKey(publisher)
	WriteJSON(output, publisher)
	radio, err := api.Register(
		"waves",
		spec.NewUser("radio@email.com", "", "", nil, "radio", "", "www.radio_station.com", "Organization"),
	)
	if err != nil {
		t.Fatal(err)
	}
	radioId := GetUserId(radio)
	radioPrivkey := GetPrivateKey(radio)
	WriteJSON(output, radio)
	if err := api.Login(composerPrivkey.String(), composerId); err != nil {
		t.Fatal(err)
	}
	composition, err := spec.NewComposition([]string{composerId}, "T-034.524.680-1", "EN", "composition_title", publisherId, "www.composition_url.com")
	if err != nil {
		t.Fatal(err)
	}
	compositionId, err := api.Publish(composition, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, Data{"compositionId": compositionId})
	sig, err := ld.ProveComposer(CHALLENGE, composerId, compositionId, composerPrivkey)
	if err != nil {
		t.Fatal(err)
	}
	if err = ld.VerifyComposer(CHALLENGE, composerId, compositionId, sig); err != nil {
		t.Fatal(err)
	}
	SleepSeconds(2)
	compositionRightId, err := api.Right(20, "", publisherId, compositionId)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, Data{"compositionRightId": compositionRightId})
	sig, err = ld.ProveRightHolder(CHALLENGE, composerPrivkey, composerId, compositionRightId)
	if err != nil {
		t.Fatal(err)
	}
	if err = ld.VerifyRightHolder(CHALLENGE, composerId, compositionRightId, sig); err != nil {
		t.Fatal(err)
	}
	sig, err = ld.ProveRightHolder(CHALLENGE, publisherPrivkey, publisherId, compositionRightId)
	if err != nil {
		t.Fatal(err)
	}
	if err = ld.VerifyRightHolder(CHALLENGE, publisherId, compositionRightId, sig); err != nil {
		t.Fatal(err)
	}
	if err = api.Login(publisherPrivkey.String(), publisherId); err != nil {
		t.Fatal(err)
	}
	mechanicalLicense, err := spec.NewLicense([]string{compositionId}, []string{performerId, producerId, recordLabelId}, publisherId, []string{compositionRightId}, "2020-01-01", "2024-01-01")
	if err != nil {
		t.Fatal(err)
	}
	mechanicalLicenseId, err := api.License(mechanicalLicense)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, Data{"mechanicalLicenseId": mechanicalLicenseId})
	sig, err = ld.ProveLicenseHolder(CHALLENGE, performerId, mechanicalLicenseId, performerPrivkey)
	if err != nil {
		t.Fatal(err)
	}
	if err = ld.VerifyLicenseHolder(CHALLENGE, performerId, mechanicalLicenseId, sig); err != nil {
		t.Fatal(err)
	}
	sig, err = ld.ProveLicenseHolder(CHALLENGE, producerId, mechanicalLicenseId, producerPrivkey)
	if err != nil {
		t.Fatal(err)
	}
	if err = ld.VerifyLicenseHolder(CHALLENGE, producerId, mechanicalLicenseId, sig); err != nil {
		t.Fatal(err)
	}
	sig, err = ld.ProveLicenseHolder(CHALLENGE, recordLabelId, mechanicalLicenseId, recordLabelPrivkey)
	if err != nil {
		t.Fatal(err)
	}
	if err = ld.VerifyLicenseHolder(CHALLENGE, recordLabelId, mechanicalLicenseId, sig); err != nil {
		t.Fatal(err)
	}

	if err = api.Login(performerPrivkey.String(), performerId); err != nil {
		t.Fatal(err)
	}
	recording, err := spec.NewRecording([]string{performerId, producerId}, compositionId, "PT2M43S", "US-S1Z-99-00001", mechanicalLicenseId, recordLabelId, "www.recording_url.com")
	if err != nil {
		t.Fatal(err)
	}
	perfomerSignature, err := api.SignRecording(recording, performerId, []int{80, 20})
	if err != nil {
		t.Fatal(err)
	}
	if err = api.Login(producerPrivkey.String(), producerId); err != nil {
		t.Fatal(err)
	}
	producerSignature, err := api.SignRecording(recording, performerId, []int{80, 20})
	if err != nil {
		t.Fatal(err)
	}
	if err = api.Login(performerPrivkey.String(), performerId); err != nil {
		t.Fatal(err)
	}
	recordingId, err := api.Release(recording, []string{perfomerSignature, producerSignature}, []int{80, 20})
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, Data{"recordingId": recordingId})
	SleepSeconds(2)
	sig, err = ld.ProveArtist(performerId, CHALLENGE, performerPrivkey, recordingId)
	if err != nil {
		t.Fatal(err)
	}
	if err = ld.VerifyArtist(performerId, CHALLENGE, recordingId, sig); err != nil {
		t.Fatal(err)
	}
	sig, err = ld.ProveArtist(producerId, CHALLENGE, producerPrivkey, recordingId)
	if err != nil {
		t.Fatal(err)
	}
	if err = ld.VerifyArtist(producerId, CHALLENGE, recordingId, sig); err != nil {
		t.Fatal(err)
	}
	recordingRightId, err := api.Right(20, "", recordLabelId, recordingId)
	WriteJSON(output, Data{"recordingRightId": recordingRightId})
	sig, err = ld.ProveRightHolder(CHALLENGE, performerPrivkey, performerId, recordingRightId)
	if err != nil {
		t.Fatal(err)
	}
	if err = ld.VerifyRightHolder(CHALLENGE, performerId, recordingRightId, sig); err != nil {
		t.Fatal(err)
	}
	sig, err = ld.ProveRightHolder(CHALLENGE, recordLabelPrivkey, recordLabelId, recordingRightId)
	if err != nil {
		t.Fatal(err)
	}
	if err = ld.VerifyRightHolder(CHALLENGE, recordLabelId, recordingRightId, sig); err != nil {
		t.Fatal(err)
	}
	if err = api.Login(recordLabelPrivkey.String(), recordLabelId); err != nil {
		t.Fatal(err)
	}
	masterLicense, err := spec.NewLicense([]string{recordingId}, []string{radioId}, recordLabelId, []string{recordingRightId}, "2020-01-01", "2022-01-01")
	if err != nil {
		t.Fatal(err)
	}
	masterLicenseId, err := api.License(masterLicense)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, Data{"masterLicenseId": masterLicenseId})
	sig, err = ld.ProveLicenseHolder(CHALLENGE, radioId, masterLicenseId, radioPrivkey)
	if err != nil {
		t.Fatal(err)
	}
	if err = ld.VerifyLicenseHolder(CHALLENGE, radioId, masterLicenseId, sig); err != nil {
		t.Fatal(err)
	}
	txs, err := bigchain.HttpGetFilter(func(txId string) (Data, error) {
		return ld.ValidateCompositionId(txId)
	}, composerPrivkey.Public())
	if err != nil {
		t.Fatal(err)
	}
	PrintJSON(txs)
}
