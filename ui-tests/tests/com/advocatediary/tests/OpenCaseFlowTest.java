package com.advocatediary.tests;

import com.advocatediary.pages.CasePage;
import com.advocatediary.pages.ClientPage;
import org.testng.Assert;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.Test;

public class OpenCaseFlowTest extends BaseTest {

  @Override
  @BeforeMethod
  public void openBrowser() {
    System.setProperty("baseUrl", "http://localhost:8080");
    System.setProperty("headless", "false");
    super.openBrowser();
  }

  @Test
  public void afterSaveClientOpenCaseThenSaveCaseDetails() {
    ClientPage client = signInToDiary()
        .nav()
        .openClients()
        .addClient()
        .saveClientWithNotes("Meera Rao", "9876543210", "meera@example.test", "Chennai", "Family matter");

    Assert.assertEquals(client.heading(), "Meera Rao");
    Assert.assertTrue(client.showsOpenACase(), "Open a case form should show after save client");

    client.fillOpenCase(
        "Rao v State",
        "Madras High Court",
        "civil",
        "WP 12 of 2026",
        "State of Tamil Nadu");
    CasePage matter = client.clickOpenCase();

    Assert.assertEquals(matter.heading(), "Rao v State");
    Assert.assertTrue(matter.showsCaseDetails(), "Case details and Save case should show");

    matter.fillCaseDetails(
        "Rao v State updated",
        "Madras High Court",
        "WP 12 of 2026",
        "writ",
        "active",
        "State of Tamil Nadu",
        "2026-01-15",
        "Listed for admission");
    matter.clickSaveCaseExpectingTitle("Rao v State updated");

    Assert.assertEquals(matter.heading(), "Rao v State updated");
    Assert.assertEquals(matter.titleValue(), "Rao v State updated");
    Assert.assertEquals(matter.notesValue(), "Listed for admission");
  }
}
