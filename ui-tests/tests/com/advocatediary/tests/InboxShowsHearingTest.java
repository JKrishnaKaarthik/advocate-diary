package com.advocatediary.tests;

import com.advocatediary.pages.CasePage;
import com.advocatediary.pages.ClientPage;
import com.advocatediary.pages.InboxPage;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;
import org.testng.Assert;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.Test;

public class InboxShowsHearingTest extends BaseTest {

  @Override
  @BeforeMethod
  public void openBrowser() {
    System.setProperty("baseUrl", "http://localhost:8080");
    System.setProperty("headless", "false");
    super.openBrowser();
  }

  @Test
  public void nextHearingShowsInInbox() {
    String stamp = String.valueOf(System.currentTimeMillis());
    String clientName = "Inbox Client " + stamp;
    String title = "Inbox Matter " + stamp;
    String caseNumber = "WP " + stamp;
    String court = "Madras High Court";
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
        .saveClientWithNotes(clientName, "9876543210", "inbox@example.test", "Chennai", "Reminder matter");

    CasePage matter = client.openCase(title, court, "civil", caseNumber, "State of Tamil Nadu");
    matter.fillHearing(when, court, "arguments");
    matter.clickSaveHearing();
    matter.waitForHearing(court, "arguments");

    InboxPage inbox = matter.nav().openInbox();
    Assert.assertTrue(inbox.isOpen());
    inbox.waitForNextHearing(clientName, title, caseNumber);
    Assert.assertTrue(
        inbox.showsNextHearing(clientName, title, caseNumber),
        "next hearing should show in Inbox with client, title, and case number");
  }
}
