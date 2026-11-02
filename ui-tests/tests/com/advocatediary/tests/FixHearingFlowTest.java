package com.advocatediary.tests;

import com.advocatediary.pages.CasePage;
import com.advocatediary.pages.ClientPage;
import org.testng.Assert;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.Test;

public class FixHearingFlowTest extends BaseTest {

  @Override
  @BeforeMethod
  public void openBrowser() {
    System.setProperty("baseUrl", "http://localhost:8080");
    System.setProperty("headless", "false");
    super.openBrowser();
  }

  @Test
  public void saveHearingWithoutWhenShowsDateRequired() {
    CasePage matter = openCaseForHearing();
    Assert.assertTrue(matter.showsFixAHearing(), "Fix a hearing and Save hearing should show");

    matter.clickSaveHearing();
    Assert.assertTrue(matter.showsDateRequired());
  }

  @Test
  public void fillHearingDetailsAndClickSaveHearing() {
    CasePage matter = openCaseForHearing();
    Assert.assertTrue(matter.showsFixAHearing(), "Fix a hearing and Save hearing should show");

    matter.fillHearing("2027-06-20T11:15", "City Civil Court", "mention");
    matter.clickSaveHearing();
    matter.waitForHearing("City Civil Court", "mention");
    Assert.assertTrue(matter.showsHearing("City Civil Court", "mention"));
  }

  private CasePage openCaseForHearing() {
    ClientPage client = signInToDiary()
        .nav()
        .openClients()
        .addClient()
        .saveClientWithNotes("Meera Rao", "9876543210", "meera@example.test", "Chennai", "Family matter");
    return client.openCase(
        "Rao v State",
        "Madras High Court",
        "civil",
        "WP 12 of 2026",
        "State of Tamil Nadu");
  }
}
