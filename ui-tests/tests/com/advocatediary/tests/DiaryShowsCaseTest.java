package com.advocatediary.tests;

import com.advocatediary.pages.CasePage;
import com.advocatediary.pages.ClientPage;
import com.advocatediary.pages.DiaryPage;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;
import org.testng.Assert;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.Test;

public class DiaryShowsCaseTest extends BaseTest {

  @Override
  @BeforeMethod
  public void openBrowser() {
    System.setProperty("baseUrl", "http://localhost:8080");
    System.setProperty("headless", "false");
    super.openBrowser();
  }

  @Test
  public void savedCaseDetailsShowOnCourtDiary() {
    String stamp = String.valueOf(System.currentTimeMillis());
    String clientName = "Diary Client " + stamp;
    String title = "Diary Matter " + stamp;
    String caseNumber = "WP " + stamp;
    String court = "Madras High Court";
    String purpose = "arguments";
    String when = LocalDateTime.now(ZoneId.of("Asia/Kolkata"))
        .plusDays(2)
        .withHour(11)
        .withMinute(15)
        .withSecond(0)
        .withNano(0)
        .format(DateTimeFormatter.ofPattern("yyyy-MM-dd'T'HH:mm"));

    ClientPage client = signInToDiary()
        .nav()
        .openClients()
        .addClient()
        .saveClientWithNotes(clientName, "9876543210", "diary@example.test", "Chennai", "Listed matter");

    CasePage matter = client.openCase(title, court, "civil", caseNumber, "State of Tamil Nadu");
    matter.fillCaseDetails(
        title,
        court,
        caseNumber,
        "civil",
        "active",
        "State of Tamil Nadu",
        "2026-01-15",
        "For diary listing");
    matter.clickSaveCaseExpectingTitle(title);
    matter.fillHearing(when, court, purpose);
    matter.clickSaveHearingExpectingSaved();

    DiaryPage diary = matter.nav().openDiary();
    Assert.assertEquals(diary.heading(), "Court diary");
    diary.waitForSavedCase(clientName, title, caseNumber, court, purpose);
    Assert.assertTrue(
        diary.showsSavedCase(clientName, title, caseNumber, court, purpose),
        "saved case title, number, court, and hearing should show on the diary");
  }
}
