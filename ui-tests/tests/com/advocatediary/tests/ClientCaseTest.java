package com.advocatediary.tests;

import com.advocatediary.pages.CasePage;
import com.advocatediary.pages.ClientPage;
import com.advocatediary.pages.ClientsPage;
import com.advocatediary.pages.DiaryPage;
import com.advocatediary.pages.InboxPage;
import org.testng.Assert;
import org.testng.annotations.Test;

public class ClientCaseTest extends BaseTest {

  @Test
  public void lawyerCanAddClientAndOpenCase() {
    DiaryPage diary = signInToDiary();
    Assert.assertTrue(diary.isOpen());

    ClientsPage clients = diary.nav().openClients();
    Assert.assertTrue(clients.isClientsPage());

    ClientPage client = clients.addClient()
        .saveClient("Meera Rao", "9876543210", "meera@example.test", "Chennai");
    Assert.assertEquals(client.heading(), "Meera Rao");

    CasePage matter = client.openCase(
        "Rao v State",
        "Madras High Court",
        "civil",
        "WP 12 of 2026",
        "State of Tamil Nadu");
    Assert.assertEquals(matter.heading(), "Rao v State");

    InboxPage inbox = diary.nav().openInbox();
    Assert.assertTrue(inbox.isOpen());
  }
}
